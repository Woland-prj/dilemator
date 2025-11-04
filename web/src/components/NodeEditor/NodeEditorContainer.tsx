import { NIL as NIL_UUID, UUIDTypes } from "uuid";
import { signal, computed } from "@preact/signals";
import { Dilemma, Node, Scenario } from "./types";
import { ScenarioForm } from "./ScenarioForm";
import { ImagePreview } from "./ImagePreview";
import { NodeNavigation } from "./NodeNavigation";
import { ScenarioList } from "./ScenarioList";

export type NodeEditorProps = {
  Dilemma: Dilemma;
  Node: Node;
  IsNew: boolean;
};

export const NodeEditorContainer = ({
  Dilemma,
  Node,
  IsNew,
}: NodeEditorProps) => {
  const imageFile = signal<File | null>(null);
  const topic = signal(Dilemma.Topic);
  const scenarioName = signal(Node.Name);
  const scenarioDescription = signal(Node.Value);
  const dilemmaId = signal<UUIDTypes>(Dilemma.ID);
  const nodeId = signal<UUIDTypes>(Node.ID);
  const scenarios = signal<Scenario[]>(Node.Scenarios);
  const isSendable = computed<boolean>((): boolean => {
    return (
      topic.value !== "" &&
      scenarioName.value !== "" &&
      scenarioDescription.value !== ""
    );
  });

  const handleSubmit = async (e: SubmitEvent) => {
    e.preventDefault();
    if (!isSendable.value) return;

    const data =
      Node.ParentID === NIL_UUID
        ? {
            topic: topic.value,
            rootName: scenarioName.value,
            rootValue: scenarioDescription.value,
          }
        : {
            name: scenarioName.value,
            value: scenarioDescription.value,
          };

    console.log(data);

    const resp = await fetch("/api/dilemma", {
      method: IsNew ? "POST" : "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });

    if (resp.ok) {
      const respData: NodeEditorProps = await resp.json();
      topic.value = respData.Dilemma.Topic;
      scenarioName.value = respData.Node.Name;
      scenarioDescription.value = respData.Node.Value;
      dilemmaId.value = respData.Dilemma.ID;
      nodeId.value = respData.Node.ID;
      scenarios.value = respData.Node.Scenarios;
      IsNew = respData.IsNew;
    }
  };

  return (
    <div class="flex items-center justify-center min-h-[92vh] text-center space-y-6 animate-fade-in">
      <div class="flex min-h-[70vh] w-full items-stretch justify-center gap-x-4">
        <NodeNavigation parentId={Node.ParentID} />

        <form
          onSubmit={handleSubmit}
          class="flex items-stretch justify-between w-4/6 rounded-box bg-base-200 border-base-300 border gap-x-4 p-4"
          enctype="multipart/form-data"
        >
          <ScenarioForm
            topic={topic}
            scenarioName={scenarioName}
            scenarioDescription={scenarioDescription}
            imageFile={imageFile}
            isSendable={isSendable}
            isRoot={Node.ParentID === NIL_UUID}
          />

          <ImagePreview imageFile={imageFile} />
        </form>

        <ScenarioList
          dilemmaId={dilemmaId}
          scenarios={scenarios}
          nodeId={nodeId}
        />
      </div>
    </div>
  );
};
