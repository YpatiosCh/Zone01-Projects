import { createApp } from "../AuraJS/index.js";
import todoReducer from "./state/reducer.js";
import { InitialState } from "./state/state.js";
import TodoList from "./views/todoList.js";

function view(state, dispatch, ctx) {
  return ctx.handler(state, dispatch, ctx);
}

const app = createApp({
  root: ".todoapp-container",
  reducer: todoReducer,
  state: InitialState,
  view,
  routes: [
    { path: "/", handler: TodoList },
    { path: "/active", handler: TodoList },
    { path: "/completed", handler: TodoList },
  ],
});

// Focus the edit input after render (needs to happen after DOM update)
const observer = new MutationObserver(() => {
  const editInput = document.querySelector(".edit");
  if (editInput && document.activeElement !== editInput) {
    editInput.focus();
    editInput.setSelectionRange(editInput.value.length, editInput.value.length);
  }
});
observer.observe(document.querySelector(".todoapp-container"), {
  childList: true,
  subtree: true,
});
