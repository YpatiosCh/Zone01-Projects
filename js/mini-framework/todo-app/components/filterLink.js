import { h } from "../../AuraJS/index.js";

export function filterLink(label, href, isSelected) {
  return h("li", {},
    h("a", {
      class: isSelected ? "selected" : "",
      href,
    }, label)
  );
}
