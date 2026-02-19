/// Hierarchy Resolver Trait
/// This trait allows the Compiler to resolve hierarchy relationships at compile time.
/// It bridges the separation between the Compute Engine and the Metadata Store.
pub trait HierarchyResolver {
    /// Returns the immediate children of a member.
    /// e.g. "North America" -> ["USA", "Canada", "Mexico"]
    fn get_children(&self, dimension: &str, member: &str) -> Vec<String>;

    /// Returns the parent of a member.
    /// e.g. "USA" -> "North America"
    fn get_parent(&self, dimension: &str, member: &str) -> Option<String>;

    /// Returns all descendants (recursive children).
    fn get_descendants(&self, dimension: &str, member: &str) -> Vec<String>;
}

/// A Mock Resolver for testing and initial development.
pub struct MockHierarchyResolver;

impl HierarchyResolver for MockHierarchyResolver {
    fn get_children(&self, _dimension: &str, member: &str) -> Vec<String> {
        match member {
            "North America" => vec!["USA".to_string(), "Canada".to_string(), "Mexico".to_string()],
            "Europe" => vec!["UK".to_string(), "France".to_string(), "Germany".to_string()],
            _ => vec![],
        }
    }

    fn get_parent(&self, _dimension: &str, member: &str) -> Option<String> {
        match member {
            "USA" => Some("North America".to_string()),
            _ => None,
        }
    }

    fn get_descendants(&self, dimension: &str, member: &str) -> Vec<String> {
        self.get_children(dimension, member) // Simple mock
    }
}
