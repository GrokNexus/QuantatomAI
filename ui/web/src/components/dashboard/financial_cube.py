# financial_cube.py
# ────────────────────────────────────────────────
# Financial Multidimensional Cube Example
# For planning & reporting (time × entity × account × scenario)
# Uses pandas + numpy + matplotlib
# ────────────────────────────────────────────────

import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from datetime import datetime, timedelta

# ────────────────────────────────────────────────
# 1. Generate realistic sample financial data
# ────────────────────────────────────────────────

def generate_financial_cube(
    n_periods=24,           # e.g. 2 years monthly
    n_entities=5,           # companies / business units
    accounts=['Revenue', 'COGS', 'OPEX', 'EBITDA', 'CashFlow'],
    scenarios=['Budget', 'Forecast', 'Actual']
):
    # Time dimension: monthly from Jan 2024
    start_date = datetime(2024, 1, 1)
    dates = [start_date + timedelta(days=30*i) for i in range(n_periods)]
    date_str = [d.strftime('%Y-%m') for d in dates]

    # Entities
    entities = [f'Entity_{chr(65+i)}' for i in range(n_entities)]  # A, B, C, ...

    # Create full Cartesian product index
    index = pd.MultiIndex.from_product(
        [date_str, entities, accounts, scenarios],
        names=['Period', 'Entity', 'Account', 'Scenario']
    )

    # Base values + some realistic variation
    np.random.seed(42)
    base_values = np.random.uniform(800_000, 1_500_000, size=len(index))

    # Add seasonality (strong Q4 peak, Q1 dip)
    month = pd.to_datetime(index.get_level_values('Period') + '-01').month
    seasonality = 1 + 0.3 * np.sin(2 * np.pi * (month - 1) / 12)

    # Scenario adjustment: Actual = Budget ± noise, Forecast optimistic
    scenario_factor = np.where(index.get_level_values('Scenario') == 'Actual', 1.0,
                              np.where(index.get_level_values('Scenario') == 'Forecast', 1.12, 1.0))
    noise = np.random.normal(1.0, 0.08, size=len(index))

    df = pd.DataFrame({
        'Value': base_values * seasonality * scenario_factor * noise
    }, index=index)

    # Round to 2 decimals (money)
    df['Value'] = df['Value'].round(2)

    return df


# ────────────────────────────────────────────────
# 2. Apply recurrence to simulate cyclical component
# ────────────────────────────────────────────────

def apply_recurrence(series: pd.Series, m: float = 1.5, phi0: float = 0.0, phi1: float = 0.0) -> pd.Series:
    """
    Apply φ(t+1) = (2 − m²) φ(t) − φ(t−1) to a time series.
    Useful for adding realistic seasonal/business cycle oscillation.
    """
    coef = 2 - m**2
    result = np.zeros(len(series))
    result[0] = phi0
    if len(series) > 1:
        result[1] = phi1

    for t in range(2, len(series)):
        result[t] = coef * result[t-1] - result[t-2]

    # Scale to match original magnitude (optional)
    if series.std() > 0:
        result = result * (series.std() / np.std(result))

    return pd.Series(result, index=series.index)


# ────────────────────────────────────────────────
# 3. Example usage
# ────────────────────────────────────────────────

if __name__ == "__main__":
    # Create the cube
    cube = generate_financial_cube(
        n_periods=36,       # 3 years
        n_entities=8,
        accounts=['Revenue', 'COGS', 'OPEX', 'EBITDA', 'NetIncome'],
        scenarios=['Budget', 'Forecast', 'Actual', 'PriorYear']
    )

    print("Cube shape:", cube.shape)
    print("\nFirst few rows:")
    print(cube.head(12))

    # ── Slice example: Revenue for Entity_A, Actual scenario
    revenue_actual = cube.xs(
        ('Revenue', 'Entity_A', 'Actual'),
        level=['Account', 'Entity', 'Scenario']
    )

    # ── Apply recurrence to simulate cycle on Revenue
    revenue_with_cycle = apply_recurrence(revenue_actual['Value'], m=1.8)

    # Plot
    plt.figure(figsize=(12, 6))
    plt.plot(revenue_actual.index, revenue_actual['Value'], label='Actual Revenue (raw)', marker='o', alpha=0.6)
    plt.plot(revenue_actual.index, revenue_with_cycle, label='With cycle (m=1.8)', linewidth=2.5)
    plt.title('Revenue - Entity A - Actual vs Cyclical Simulation')
    plt.xlabel('Period')
    plt.ylabel('USD')
    plt.legend()
    plt.grid(True, alpha=0.3)
    plt.xticks(rotation=45)
    plt.tight_layout()
    plt.show()

    # ── Aggregation example: Total EBITDA by Scenario
    ebitda_by_scenario = cube.xs('EBITDA', level='Account').groupby('Scenario').sum()
    print("\nTotal EBITDA by Scenario:")
    print(ebitda_by_scenario)

    # ── Pivot example: Revenue by Entity × Period (Actual only)
    revenue_pivot = cube.xs(
        ('Revenue', 'Actual'),
        level=['Account', 'Scenario']
    ).unstack('Entity')['Value']
    print("\nRevenue Pivot (Entity × Period):")
    print(revenue_pivot.head())

    # ── Export full cube to Excel
    cube_reset = cube.reset_index()
    cube_reset.to_excel('financial_cube_export.xlsx', index=False)
    print("\nExported to financial_cube_export.xlsx")