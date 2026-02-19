use petgraph::graph::{DiGraph, NodeIndex};
use petgraph::algo::toposort;
use std::collections::HashMap;

/// The Dependency Graph tracks relationships between Atoms/Dimensions.
/// e.g. "Net Income" -> "Tax" -> "Revenue"
pub struct DependencyGraph {
    graph: DiGraph<String, ()>,
    node_map: HashMap<String, NodeIndex>,
}

impl DependencyGraph {
    pub fn new() -> Self {
        Self {
            graph: DiGraph::new(),
            node_map: HashMap::new(),
        }
    }

    /// Adds a node (e.g., "Revenue") to the graph if it doesn't exist.
    pub fn add_node(&mut self, name: &str) -> NodeIndex {
        if let Some(&idx) = self.node_map.get(name) {
            return idx;
        }
        let idx = self.graph.add_node(name.to_string());
        self.node_map.insert(name.to_string(), idx);
        idx
    }

    /// Adds a dependency: `from` depends on `to`.
    /// e.g., add_dependency("Net Income", "Revenue") means Revenue must be calc'd first.
    /// In graph terms: Edge from Revenue -> Net Income.
    pub fn add_dependency(&mut self, dependent: &str, dependency: &str) {
        let u = self.add_node(dependency); // Revenue
        let v = self.add_node(dependent);  // Net Income
        self.graph.add_edge(u, v, ());
    }

    /// Returns the execution order (Topological Sort).
    /// Items at the start of the list should be calculated first.
    pub fn resolve_order(&self) -> Result<Vec<String>, String> {
        match toposort(&self.graph, None) {
            Ok(nodes) => {
                let order: Vec<String> = nodes
                    .iter()
                    .map(|&idx| self.graph[idx].clone())
                    .collect();
                Ok(order)
            }
            Err(_) => Err("Cycle detected in dependency graph!".to_string()),
        }
    }
}
