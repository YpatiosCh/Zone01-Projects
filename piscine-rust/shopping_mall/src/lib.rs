pub mod mall;

pub use mall::*;
use std::collections::HashMap;

pub fn biggest_store(m: &Mall) -> (&str, &Store) {
    m.floors
        .values()
        .flat_map(|f| f.stores.iter())
        .max_by_key(|(_, s)| s.square_meters)
        .map(|(name, store)| (name.as_str(), store))
        .unwrap()
}

pub fn highest_paid_employee(m: &Mall) -> Vec<(&str, &Employee)> {
    let all: Vec<(&str, &Employee)> = m
        .floors
        .values()
        .flat_map(|f| f.stores.values())
        .flat_map(|s| s.employees.iter())
        .map(|(name, emp)| (name.as_str(), emp))
        .collect();

    let max_salary = all
        .iter()
        .map(|(_, e)| e.salary)
        .fold(f64::NEG_INFINITY, f64::max);

    all.into_iter()
        .filter(|(_, e)| e.salary == max_salary)
        .collect()
}

pub fn nbr_of_employees(m: &Mall) -> usize {
    let employees: usize = m
        .floors
        .values()
        .flat_map(|f| f.stores.values())
        .map(|s| s.employees.len())
        .sum();

    employees + m.guards.len()
}

pub fn check_for_securities(mall: &mut Mall, available_sec: HashMap<String, Guard>) {
    let total_size = mall.floors.values().map(|f| f.size_limit).sum::<u64>();

    let total_areas = total_size / 200;
    let unguarded_areas = total_areas as usize - mall.guards.len();

    available_sec
        .into_iter()
        .take(unguarded_areas)
        .for_each(|(name, guard)| {
            mall.hire_guard(name, guard);
        });
}

pub fn cut_or_raise(m: &mut Mall) {
    m.floors
        .values_mut()
        .flat_map(|f| f.stores.values_mut())
        .flat_map(|s| s.employees.values_mut())
        .for_each(|e| {
            let hours = e.working_hours.1 - e.working_hours.0;
            if hours >= 10 {
                e.raise(e.salary * 0.1);
            } else {
                e.cut(e.salary * 0.1);
            }
        });
}

#[cfg(test)]
mod tests {
    use super::*;

    fn make_mall() -> Mall {
        Mall::new(
            "Test Mall",
            [(
                "Guard One",
                Guard {
                    age: 30,
                    years_experience: 5,
                },
            )]
            .into(),
            [(
                "Floor A",
                Floor::new(
                    [
                        (
                            "Big Store",
                            Store::new(
                                [(
                                    "Alice",
                                    Employee {
                                        age: 25,
                                        working_hours: (8, 20),
                                        salary: 500.0,
                                    },
                                )]
                                .into(),
                                500,
                            ),
                        ),
                        (
                            "Small Store",
                            Store::new(
                                [(
                                    "Bob",
                                    Employee {
                                        age: 30,
                                        working_hours: (9, 14),
                                        salary: 1200.0,
                                    },
                                )]
                                .into(),
                                100,
                            ),
                        ),
                    ]
                    .into(),
                    1000,
                ),
            )]
            .into(),
        )
    }

    #[test]
    fn test_biggest_store() {
        let mall = make_mall();
        let (name, store) = biggest_store(&mall);
        assert_eq!(name, "Big Store");
        assert_eq!(store.square_meters, 500);
    }

    #[test]
    fn test_highest_paid_employee() {
        let mall = make_mall();
        let result = highest_paid_employee(&mall);
        assert_eq!(result.len(), 1);
        assert_eq!(result[0].0, "Bob");
        assert_eq!(result[0].1.salary, 1200.0);
    }

    #[test]
    fn test_nbr_of_employees() {
        let mall = make_mall();

        assert_eq!(nbr_of_employees(&mall), 3);
    }

    #[test]
    fn test_cut_or_raise() {
        let mut mall = make_mall();
        cut_or_raise(&mut mall);
        let floor = mall.floors.get("Floor A").unwrap();

        assert_eq!(floor.stores["Big Store"].employees["Alice"].salary, 550.0);

        assert_eq!(floor.stores["Small Store"].employees["Bob"].salary, 1080.0);
    }
}
