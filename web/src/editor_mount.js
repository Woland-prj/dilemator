import { h, render } from "preact";
import NodeEditor from "./components/NodeEditor.jsx";

// Монтируем компонент в контейнер, если он есть на странице
const root = document.getElementById("editor-root");
if (root) {
    const props = JSON.parse(root.getAttribute('data-props')) || { Topic: "New dilemma", Node: {} };
    render(
        <NodeEditor props={props} />,
        root,
    );
}