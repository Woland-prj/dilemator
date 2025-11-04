import { NIL as NIL_UUID, UUIDTypes } from "uuid";
import { Scenario } from "./types";
import { Signal } from "@preact/signals";

type ScenarioListProps = {
  dilemmaId: Signal<UUIDTypes>;
  scenarios: Signal<Scenario[]>;
  nodeId: Signal<UUIDTypes>;
};

export const ScenarioList = ({
  dilemmaId,
  scenarios,
  nodeId,
}: ScenarioListProps) => {
  return (
    <div class="flex flex-col gap-y-4">
      {scenarios.value?.map((s) => (
        <button key={s.ID} class="btn btn-neutral btn-block">
          {s.Name}
        </button>
      ))}
      <button
        class={`btn btn-primary btn-block ${nodeId.value === NIL_UUID ? "btn-disabled" : ""}`}
        hx-get={`/components/dilemma/editor?did=${dilemmaId.value}`}
        hx-target="#content-container"
        hx-swap="innerHTML"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 4v16m8-8H4"
          ></path>
        </svg>
      </button>
    </div>
  );
};
