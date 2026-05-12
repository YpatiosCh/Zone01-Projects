#[derive(Clone, Debug)]
pub struct List<T> {
    pub head: Option<Node<T>>,
}

#[derive(Clone, Debug)]
pub struct Node<T> {
    pub value: T,
    pub next: Option<Box<Node<T>>>,
}

impl<T> List<T> {
    pub fn new() -> Self {
        List { head: None }
    }

    pub fn push(&mut self, value: T) {
        let mut new_node = Node { value, next: None };

        match self.head.take() {
            None => self.head = Some(new_node),
            Some(head) => {
                new_node.next = Some(Box::new(head));
                self.head = Some(new_node);
            }
        }
    }

    pub fn pop(&mut self) {
        if let Some(node) = self.head.take() {
            self.head = node.next.map(|n| *n);
        }
    }

    pub fn len(&self) -> usize {
        let mut len: usize = 0;
        let mut current_node = self.head.as_ref();

        while let Some(node) = current_node {
            len += 1;
            current_node = node.next.as_deref();
        }

        len
    }
}
