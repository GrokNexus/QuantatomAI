import json
import importlib
import logging
import os
import time
import uuid
from typing import Any, Dict, List, Optional

logger = logging.getLogger("cortex.db")


def _get_database_url() -> str:
    return os.getenv("DATABASE_URL", "").strip()


def _is_valid_uuid(value: str) -> bool:
    try:
        uuid.UUID(value)
        return True
    except Exception:
        return False


def log_ai_inference_event(
    *,
    tenant_id: str,
    request_type: str,
    model_provider: str,
    model_id: str,
    confidence_score: float,
    request_payload: Dict[str, Any],
    response_payload: Dict[str, Any],
    grounding_atoms: Optional[List[str]] = None,
    app_id: Optional[str] = None,
    human_override: bool = False,
    override_reason: Optional[str] = None,
    inference_latency_ms: Optional[int] = None,
) -> None:
    """
    Best-effort persistence for ai_inference_log.
    This function intentionally never raises to avoid breaking inference API paths.
    """
    db_url = _get_database_url()
    if not db_url:
        return

    try:
        psycopg = importlib.import_module("psycopg")
    except Exception:
        logger.warning("psycopg not installed; skipping ai_inference_log write")
        return

    if not _is_valid_uuid(tenant_id):
        logger.warning("invalid tenant_id for ai inference log: %s", tenant_id)
        return

    if app_id and not _is_valid_uuid(app_id):
        logger.warning("invalid app_id for ai inference log: %s", app_id)
        app_id = None

    grounding = grounding_atoms or []

    started = time.time()
    try:
        with psycopg.connect(db_url) as conn:
            with conn.cursor() as cur:
                cur.execute(
                    """
                    INSERT INTO ai_inference_log (
                        tenant_id,
                        app_id,
                        request_type,
                        model_provider,
                        model_id,
                        confidence_score,
                        request_payload,
                        response_payload,
                        grounding_atoms,
                        human_override,
                        override_reason,
                        inference_latency_ms
                    ) VALUES (
                        %(tenant_id)s,
                        %(app_id)s,
                        %(request_type)s,
                        %(model_provider)s,
                        %(model_id)s,
                        %(confidence_score)s,
                        %(request_payload)s::jsonb,
                        %(response_payload)s::jsonb,
                        %(grounding_atoms)s::jsonb,
                        %(human_override)s,
                        %(override_reason)s,
                        %(inference_latency_ms)s
                    )
                    """,
                    {
                        "tenant_id": tenant_id,
                        "app_id": app_id,
                        "request_type": request_type,
                        "model_provider": model_provider,
                        "model_id": model_id,
                        "confidence_score": confidence_score,
                        "request_payload": json.dumps(request_payload),
                        "response_payload": json.dumps(response_payload),
                        "grounding_atoms": json.dumps(grounding),
                        "human_override": human_override,
                        "override_reason": override_reason,
                        "inference_latency_ms": inference_latency_ms,
                    },
                )
            conn.commit()
    except Exception as ex:
        logger.warning("failed to persist ai_inference_log event: %s", ex)
        return

    elapsed_ms = int((time.time() - started) * 1000)
    logger.debug("ai inference log persisted in %d ms", elapsed_ms)
