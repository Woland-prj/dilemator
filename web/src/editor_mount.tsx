import { render } from "preact";
import {
  NodeEditorContainer,
  NodeEditorProps,
} from "./components/NodeEditor/NodeEditorContainer";
import htmx from "htmx.org";

const root = document.getElementById("editor-root");
if (root) {
  const props: NodeEditorProps = JSON.parse(
    root.getAttribute("data-props") || "",
  );
  console.log(props);
  render(<NodeEditorContainer {...props} />, root);
  htmx.process(root);
}
