use std::cell::{Cell, RefCell};

#[derive(Debug)]
pub struct ThreadPool {
    pub drops: Cell<usize>,
    pub states: RefCell<Vec<bool>>,
}

impl ThreadPool {
    pub fn new() -> Self {
        ThreadPool {
            drops: Default::default(),
            states: Default::default(),
        }
    }

    pub fn new_thread(&self, c: String) -> (usize, Thread<'_>) {
        let pid = self.thread_len();
        self.states.borrow_mut().push(false);
        (pid, Thread::new(pid, c, self))
    }

    pub fn thread_len(&self) -> usize {
        self.states.borrow().len()
    }

    pub fn is_dropped(&self, id: usize) -> bool {
        self.states.borrow()[id]
    }

    pub fn drop_thread(&self, id: usize) {
        if self.is_dropped(id) {
            panic!("{} is already dropped", id);
        }
        self.states.borrow_mut()[id] = true;
        self.drops.set(self.drops.get() + 1);
    }
}

#[derive(Debug)]
pub struct Thread<'a> {
    pub pid: usize,
    pub cmd: String,
    pub parent: &'a ThreadPool,
}

impl<'a> Thread<'a> {
    pub fn new(pid: usize, cmd: String, tp: &'a ThreadPool) -> Self {
        Thread {
            pid,
            cmd,
            parent: tp,
        }
    }

    pub fn skill(self) {
        // taking ownership means Drop fires at end of this scope
    }
}

impl Drop for Thread<'_> {
    fn drop(&mut self) {
        self.parent.drop_thread(self.pid);
    }
}
