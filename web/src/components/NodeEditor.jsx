import { h } from "preact";
import { useState } from "preact/hooks";

export default function NodeEditor(props) {
  const [imagePreview, setImagePreview] = useState(null);
  const [imageFile, setImageFile] = useState(null);
  const [scenarioName, setScenarioName] = useState("");
  const [scenarioDescription, setScenarioDescription] = useState("");

  const handleFileChange = (e) => {
    const file = e.target.files?.[0];
    if (file) {
      setImageFile(file);
      const reader = new FileReader();
      reader.onloadend = () => setImagePreview(reader.result);
      reader.readAsDataURL(file);
    } else {
      setImagePreview(null);
      setImageFile(null);
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append("name", scenarioName);
    formData.append("description", scenarioDescription);
    if (imageFile) formData.append("image", imageFile);

    fetch("/api/node/upload", {
      method: "POST",
      body: formData,
    })
      .then((res) => res.json())
      .then((data) => {
        console.log("Uploaded:", data);
      })
      .catch((err) => console.error(err));
  };

  return (
    <div class="flex items-center justify-center min-h-[92vh] text-center space-y-6 animate-fade-in">
      <div class="flex min-h-[70vh] w-full items-stretch justify-center gap-x-4">
        {props.node?.ParentID ? (
          <button class="btn btn-outline btn-dash h-auto px-1">
            <svg viewBox="0 0 24 24" fill="currentColor" class="size-8">
              <path d="M14.2893 5.70708C13.8988 5.31655 13.2657 5.31655 12.8751 5.70708L7.98768 10.5993C7.20729 11.3805 7.2076 12.6463 7.98837 13.427L12.8787 18.3174C13.2693 18.7079 13.9024 18.7079 14.293 18.3174C14.6835 17.9269 14.6835 17.2937 14.293 16.9032L10.1073 12.7175C9.71678 12.327 9.71678 11.6939 10.1073 11.3033L14.2893 7.12129C14.6799 6.73077 14.6799 6.0976 14.2893 5.70708Z"></path>
            </svg>
          </button>
        ) : (
          <button class="btn btn-outline btn-dash btn-disabled h-auto px-1">
            <svg viewBox="0 0 24 24" fill="currentColor" class="size-8">
              <path d="M14.2893 5.70708C13.8988 5.31655 13.2657 5.31655 12.8751 5.70708L7.98768 10.5993C7.20729 11.3805 7.2076 12.6463 7.98837 13.427L12.8787 18.3174C13.2693 18.7079 13.9024 18.7079 14.293 18.3174C14.6835 17.9269 14.6835 17.2937 14.293 16.9032L10.1073 12.7175C9.71678 12.327 9.71678 11.6939 10.1073 11.3033L14.2893 7.12129C14.6799 6.73077 14.6799 6.0976 14.2893 5.70708Z"></path>
            </svg>
          </button>
        )}

        <form
          onSubmit={handleSubmit}
          class="flex items-stretch justify-between w-4/6 rounded-box bg-base-200 border-base-300 border gap-x-4 p-4"
          enctype="multipart/form-data"
        >
          <fieldset class="fieldset bg-base-300 border-base-300 rounded-box w-1/2 border p-4">
            <legend class="fieldset-legend">Scenario editing</legend>
            <label class="label">Scenario name</label>
            <input
              type="text"
              value={scenarioName}
              onInput={(e) => setScenarioName(e.target.value)}
              class="input w-full"
              placeholder="Name of your scenario"
            />
            <label class="label">Scenario description</label>
            <textarea
              value={scenarioDescription}
              onInput={(e) => setScenarioDescription(e.target.value)}
              class="textarea w-full h-60"
              placeholder="Describe what's happening..."
            ></textarea>
            <label class="label">Scenario image</label>
            <input
              type="file"
              class="file-input w-full"
              onChange={handleFileChange}
            />
            <button type="submit" class="btn btn-neutral mt-4">
              Save scenario
            </button>
          </fieldset>

          <div class="flex items-center justify-center bg-base-300 border-neutral rounded-box w-1/2 border border-dashed p-4">
            {imagePreview ? (
              <img
                src={imagePreview}
                alt="Preview"
                class="max-h-96 rounded-lg shadow-lg object-contain"
              />
            ) : (
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="0.9"
                class="size-30 text-base-content/40"
              >
                <path
                  d="M14.2639 15.9375L12.5958 14.2834C11.7909 13.4851 11.3884 13.086 10.9266 12.9401C10.5204 12.8118 10.0838 12.8165 9.68048 12.9536C9.22188 13.1095 8.82814 13.5172 8.04068 14.3326L4.04409 18.2801M14.2639 15.9375L14.6053 15.599C15.4112 14.7998 15.8141 14.4002 16.2765 14.2543C16.6831 14.126 17.12 14.1311 17.5236 14.2687C17.9824 14.4251 18.3761 14.8339 19.1634 15.6514L20 16.4934M14.2639 15.9375L18.275 19.9565M18.275 19.9565C17.9176 20 17.4543 20 16.8 20H7.2C6.07989 20 5.51984 20 5.09202 19.782C4.71569 19.5903 4.40973 19.2843 4.21799 18.908C4.12796 18.7313 4.07512 18.5321 4.04409 18.2801M18.275 19.9565C18.5293 19.9256 18.7301 19.8727 18.908 19.782C19.2843 19.5903 19.5903 19.2843 19.782 18.908C20 18.4802 20 17.9201 20 16.8V16.4934M4.04409 18.2801C4 17.9221 4 17.4575 4 16.8V7.2C4 6.0799 4 5.51984 4.21799 5.09202C4.40973 4.71569 4.71569 4.40973 5.09202 4.21799C5.51984 4 6.07989 4 7.2 4H16.8C17.9201 4 18.4802 4 18.908 4.21799C19.2843 4.40973 19.5903 4.71569 19.782 5.09202C20 5.51984 20 6.0799 20 7.2V16.4934M17 8.99989C17 10.1045 16.1046 10.9999 15 10.9999C13.8954 10.9999 13 10.1045 13 8.99989C13 7.89532 13.8954 6.99989 15 6.99989C16.1046 6.99989 17 7.89532 17 8.99989Z"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                ></path>
              </svg>
            )}
          </div>
        </form>

        <div class="flex flex-col gap-y-4">
          {props.node?.Scenarios?.map((s) => (
            <button key={s.id} class="btn btn-neutral btn-block">
              {s.name}
            </button>
          ))}
          <button class="btn btn-primary btn-block">
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
      </div>
    </div>
  );
}
