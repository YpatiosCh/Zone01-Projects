use json::object;

pub struct Food {
    pub name: String,
    pub calories: (String, String),
    pub fats: f64,
    pub carbs: f64,
    pub proteins: f64,
    pub nbr_of_portions: f64,
}

pub fn calculate_macros(foods: &[Food]) -> json::JsonValue {
    let mut cals = 0.0_f64;
    let mut carbs = 0.0_f64;
    let mut proteins = 0.0_f64;
    let mut fats = 0.0_f64;

    for food in foods {
        let p = food.nbr_of_portions;

        let kcal: f64 = food
            .calories
            .1
            .trim_end_matches("kcal")
            .parse()
            .unwrap_or(0.0);

        cals += kcal * p;
        carbs += food.carbs * p;
        proteins += food.proteins * p;
        fats += food.fats * p;
    }

    object! {
        "cals":     round(cals),
        "carbs":    round(carbs),
        "proteins": round(proteins),
        "fats":     round(fats),
    }
}

fn round(val: f64) -> f64 {
    (val * 100.0).round() / 100.0
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_single_food_one_portion() {
        let foods = [Food {
            name: "chicken".to_owned(),
            calories: ("0kJ".to_owned(), "200kcal".to_owned()),
            proteins: 30.0,
            fats: 5.0,
            carbs: 0.0,
            nbr_of_portions: 1.0,
        }];
        let result = calculate_macros(&foods);
        assert_eq!(result["cals"], 200.0);
        assert_eq!(result["proteins"], 30.0);
        assert_eq!(result["fats"], 5.0);
        assert_eq!(result["carbs"], 0.0);
    }

    #[test]
    fn test_portions_multiply_macros() {
        let foods = [Food {
            name: "rice".to_owned(),
            calories: ("0kJ".to_owned(), "100kcal".to_owned()),
            proteins: 2.0,
            fats: 0.5,
            carbs: 20.0,
            nbr_of_portions: 3.0,
        }];
        let result = calculate_macros(&foods);
        assert_eq!(result["cals"], 300.0);
        assert_eq!(result["proteins"], 6.0);
        assert_eq!(result["fats"], 1.5);
        assert_eq!(result["carbs"], 60.0);
    }
}
