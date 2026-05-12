import { h } from "../../AuraJS/index.js";
import { todoItem } from "../components/todoItem.js";
import { filterLink } from "../components/filterLink.js";

export default function TodoList(state, dispatch, ctx) {
  const filter =
    ctx.path === "/active"
      ? "active"
      : ctx.path === "/completed"
        ? "completed"
        : "all";

  const filteredTodos = state.todos.filter((t) => {
    if (filter === "active") return !t.completed;
    if (filter === "completed") return t.completed;
    return true;
  });

  const activeTodoCount = state.todos.filter((t) => !t.completed).length;

  return h("section", { class: "todoapp" },
    // Header
    h("header", { class: "header" },
      h("h1", {}, "todos"),
      h("input", {
        class: "new-todo",
        placeholder: "What needs to be done?",
        autofocus: true,
        value: state.newTodo,
        onInput: (e) =>
          dispatch({ type: "SET_NEW_TODO", text: e.target.value }),
        onKeydown: (e) => {
          if (e.key === "Enter" && e.target.value.trim()) {
            dispatch({ type: "ADD_TODO", text: e.target.value.trim() });
          }
        },
      })
    ),

    // Main section (only if there are todos)
    ...(state.todos.length > 0
      ? [
        h("section", { class: "main" },
          h("input", {
            id: "toggle-all",
            class: "toggle-all",
            type: "checkbox",
            checked: activeTodoCount === 0,
            onChange: () => dispatch({ type: "TOGGLE_ALL" }),
          }),
          h("label", { for: "toggle-all" }, "Mark all as complete"),
          h("ul", { class: "todo-list" },
            ...filteredTodos.map((todo) => todoItem(todo, state, dispatch))
          )
        ),

        // Footer
        h("footer", { class: "footer" },
          h("span", { class: "todo-count" },
            h("strong", {}, String(activeTodoCount)),
            activeTodoCount === 1 ? " item left!" : " items left!"
          ),
          h("ul", { class: "filters" },
            filterLink("All", "#/", filter === "all"),
            filterLink("Active", "#/active", filter === "active"),
            filterLink("Completed", "#/completed", filter === "completed")
          ),
          h("button", {
            class: "clear-completed",
            onClick: () => dispatch({ type: "CLEAR_COMPLETED" }),
          }, "Clear completed"),

        ),
      ]
      : [])
  );
}
