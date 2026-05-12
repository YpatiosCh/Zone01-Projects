export default function todoReducer(state, action) {
  switch (action.type) {
    case "ADD_TODO": {
      const maxId = state.todos.reduce((max, t) => Math.max(max, t.id), 0);
      return {
        ...state,
        todos: [
          { id: maxId + 1, text: action.text, completed: false },
          ...state.todos,
        ],
        newTodo: "",
      };
    }

    case "TOGGLE_TODO":
      return {
        ...state,
        todos: state.todos.map((t) =>
          t.id === action.id ? { ...t, completed: !t.completed } : t
        ),
      };

    case "DELETE_TODO":
      return {
        ...state,
        todos: state.todos.filter((t) => t.id !== action.id),
      };

    case "EDIT_TODO":
      return {
        ...state,
        editing: action.id,
      };

    case "SAVE_TODO":
      return {
        ...state,
        todos: state.todos.map((t) =>
          t.id === action.id ? { ...t, text: action.text } : t
        ),
        editing: null,
      };

    case "CANCEL_EDIT":
      return {
        ...state,
        editing: null,
      };

    case "TOGGLE_ALL": {
      const allCompleted = state.todos.every((t) => t.completed);
      return {
        ...state,
        todos: state.todos.map((t) => ({ ...t, completed: !allCompleted })),
      };
    }

    case "CLEAR_COMPLETED":
      return {
        ...state,
        todos: state.todos.filter((t) => !t.completed),
      };

    case "SET_NEW_TODO":
      return {
        ...state,
        newTodo: action.text,
      };

    default:
      return state;
  }
}
