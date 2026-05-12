import { h } from "../../AuraJS/index.js";

let editCommitted = false;

export function todoItem(todo, state, dispatch) {
  const isEditing = state.editing === todo.id;
  const classes = [
    todo.completed ? "completed" : "",
    isEditing ? "editing" : "",
  ]
    .filter(Boolean)
    .join(" ");

  return h("li", { class: classes, "data-id": String(todo.id) },
    h("div", { class: "view" },
      h("input", {
        class: "toggle",
        type: "checkbox",
        checked: todo.completed,
        onChange: () => dispatch({ type: "TOGGLE_TODO", id: todo.id }),
      }),
      h("label", {
        onDblclick: () => dispatch({ type: "EDIT_TODO", id: todo.id }),
      }, todo.text),
      h("button", {
        class: "destroy",
        onClick: () => dispatch({ type: "DELETE_TODO", id: todo.id }),
      })
    ),
    ...(isEditing
      ? [
          h("input", {
            class: "edit",
            value: todo.text,
            onKeydown: (e) => {
              if (e.key === "Enter") {
                editCommitted = true;
                const text = e.target.value.trim();
                if (text) {
                  dispatch({ type: "SAVE_TODO", id: todo.id, text });
                } else {
                  dispatch({ type: "DELETE_TODO", id: todo.id });
                }
              } else if (e.key === "Escape") {
                editCommitted = true;
                dispatch({ type: "CANCEL_EDIT" });
              }
            },
            onBlur: (e) => {
              if (editCommitted) {
                editCommitted = false;
                return;
              }
              const text = e.target.value.trim();
              if (text) {
                dispatch({ type: "SAVE_TODO", id: todo.id, text });
              } else {
                dispatch({ type: "DELETE_TODO", id: todo.id });
              }
            },
          }),
        ]
      : [])
  );
}
