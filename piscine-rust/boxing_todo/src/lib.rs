mod err;

pub use err::{ParseErr, ReadErr};
use std::{error::Error, fs};

#[derive(Debug, Eq, PartialEq)]
pub struct Task {
    pub id: u32,
    pub description: String,
    pub level: u32,
}

#[derive(Debug, Eq, PartialEq)]
pub struct TodoList {
    pub title: String,
    pub tasks: Vec<Task>,
}

impl TodoList {
    pub fn get_todo(path: &str) -> Result<TodoList, Box<dyn Error>> {
        let content = fs::read_to_string(path).map_err(|e| ReadErr {
            child_err: Box::new(e),
        })?;

        let parsed = json::parse(&content).map_err(|e| ParseErr::Malformed(Box::new(e)))?;

        if parsed["tasks"].is_empty() {
            return Err(Box::new(ParseErr::Empty));
        }

        let title = parsed["title"].as_str().ok_or(ParseErr::Empty)?.to_string();

        let tasks = parsed["tasks"]
            .members()
            .map(|t| Task {
                id: t["id"].as_u32().unwrap_or(0),
                description: t["description"].as_str().unwrap_or("").to_string(),
                level: t["level"].as_u32().unwrap_or(0),
            })
            .collect();

        Ok(TodoList { title, tasks })
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::{fs::File, io::Write};

    fn create_file(name: &str, content: &str) {
        File::create(name)
            .unwrap()
            .write_all(content.as_bytes())
            .unwrap();
    }

    #[test]
    fn valid_todo() {
        create_file(
            "test_valid.json",
            r#"{
                "title": "TEST LIST",
                "tasks": [
                    { "id": 0, "description": "do this", "level": 0 },
                    { "id": 1, "description": "do that", "level": 5 }
                ]
            }"#,
        );

        let result = TodoList::get_todo("test_valid.json");
        assert!(result.is_ok());

        let list = result.unwrap();
        assert_eq!(list.title, "TEST LIST");
        assert_eq!(list.tasks.len(), 2);
        assert_eq!(
            list.tasks[0],
            Task {
                id: 0,
                description: "do this".to_string(),
                level: 0
            }
        );
        assert_eq!(
            list.tasks[1],
            Task {
                id: 1,
                description: "do that".to_string(),
                level: 5
            }
        );
    }

    #[test]
    fn empty_tasks() {
        create_file(
            "test_empty.json",
            r#"{
                "title": "TEST LIST",
                "tasks": []
            }"#,
        );

        let result = TodoList::get_todo("test_empty.json");
        assert!(result.is_err());

        let err = result.unwrap_err();
        assert_eq!(err.to_string(), "Failed to parse todo file");
        assert!(err.source().is_none()); // ParseErr::Empty → source is None
    }

    #[test]
    fn malformed_json() {
        create_file("test_malformed.json", r#"{ "something": , }"#);

        let result = TodoList::get_todo("test_malformed.json");
        assert!(result.is_err());

        let err = result.unwrap_err();
        assert_eq!(err.to_string(), "Failed to parse todo file");
        assert!(err.source().is_some()); // ParseErr::Malformed → source is Some
    }

    #[test]
    fn file_not_found() {
        let result = TodoList::get_todo("nonexistent.json");
        assert!(result.is_err());

        let err = result.unwrap_err();
        assert_eq!(err.to_string(), "Failed to read todo file");
        assert!(err.source().is_some()); // ReadErr → always has a source
    }
}
