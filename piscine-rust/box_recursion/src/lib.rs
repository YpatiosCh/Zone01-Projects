#[derive(Debug, PartialEq, Clone, Copy)]
pub enum Role {
    CEO,
    Manager,
    Worker,
}

impl From<&str> for Role {
    fn from(str_role: &str) -> Self {
        match str_role {
            "CEO" => Role::CEO,
            "Manager" => Role::Manager,
            "Normal Worker" => Role::Worker,
            _ => unreachable!(),
        }
    }
}

#[derive(Debug)]
pub struct WorkEnvironment {
    pub grade: Link,
}

pub type Link = Option<Box<Worker>>;

#[derive(Debug)]
pub struct Worker {
    pub role: Role,
    pub name: String,
    pub next: Link,
}

impl WorkEnvironment {
    pub fn new() -> Self {
        WorkEnvironment { grade: None }
    }

    pub fn add_worker(&mut self, name: &str, role: &str) {
        let worker = Worker {
            role: Role::from(role),
            name: name.to_owned(),
            next: self.grade.take(),
        };

        self.grade = Some(Box::new(worker))
    }

    pub fn remove_worker(&mut self) -> Option<String> {
        self.grade.take().map(|worker| {
            self.grade = worker.next;
            worker.name
        })
    }

    pub fn last_worker(&self) -> Option<(String, Role)> {
        self.grade
            .as_ref()
            .map(|worker| (worker.name.clone(), worker.role))
    }
}
