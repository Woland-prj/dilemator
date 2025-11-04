import { TargetedInputEvent } from "preact";
import { Signal, ReadonlySignal } from "@preact/signals";

type ScenarioFormProps = {
  topic: Signal<string>;
  scenarioName: Signal<string>;
  scenarioDescription: Signal<string>;
  imageFile: Signal<File | null>;
  isSendable: ReadonlySignal<boolean>;
  isRoot: boolean;
};

export const ScenarioForm = ({
  topic,
  scenarioName,
  scenarioDescription,
  imageFile,
  isSendable,
  isRoot,
}: ScenarioFormProps) => {
  return (
    <div class="flex flex-col h-full w-full gap-y-4">
      {isRoot ? (
        <input
          type="text"
          placeholder="Topic of your dilemma..."
          class="input w-full"
          value={topic.value}
          onInput={(e: TargetedInputEvent<HTMLInputElement>) =>
            (topic.value = e.currentTarget.value)
          }
        />
      ) : (
        <h2 class="text-2xl font-semibold text-primary-content">
          {topic.value}
        </h2>
      )}
      <fieldset class="fieldset flex flex-col bg-base-300 border-base-300 rounded-box w-full h-full border p-4">
        <legend class="fieldset-legend">Scenario editing</legend>

        <label class="label">Scenario name</label>
        <input
          type="text"
          value={scenarioName.value}
          onInput={(e: TargetedInputEvent<HTMLInputElement>) =>
            (scenarioName.value = e.currentTarget.value)
          }
          class="input w-full"
          placeholder="Name of your scenario"
        />

        <label class="label">Scenario description</label>
        <textarea
          value={scenarioDescription.value}
          onInput={(e: TargetedInputEvent<HTMLTextAreaElement>) =>
            (scenarioDescription.value = e.currentTarget.value)
          }
          class="textarea w-full flex-1 resize-none"
          placeholder="Describe what's happening..."
        ></textarea>

        <label class="label">Scenario image</label>
        <input
          type="file"
          class="file-input w-full"
          onChange={(e: TargetedInputEvent<HTMLInputElement>) =>
            (imageFile.value = e.currentTarget.files?.[0] ?? null)
          }
        />

        <button
          type="submit"
          class={`btn btn-neutral mt-4 ${isSendable.value ? "" : "btn-disabled"}`}
        >
          Save scenario
        </button>
      </fieldset>
    </div>
  );
};
