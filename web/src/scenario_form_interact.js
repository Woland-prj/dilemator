const form = document.getElementById("editor-form");
if (form) {
  const topicInput = form.querySelector('input[name="topic"]');
  const nameInput = form.querySelector('input[name="name"]');
  const descTextarea = form.querySelector('textarea[name="value"]');
  const submitButton = form.querySelector('button[type="submit"]');
  const fileInput = form.querySelector('input[name="image"]');

  if (!submitButton || !nameInput || !descTextarea) return;

  const topicInit = topicInput ? topicInput.value : undefined;
  const nameInit = nameInput.value;
  const descInit = descTextarea.value;

  const checkValidity = () => {
    const topicChanged = topicInput ? topicInput.value !== topicInit : false;
    const nameChanged = nameInput.value !== nameInit;
    const descChanged = descTextarea.value !== descInit;
    const fileChanged = fileInput ? fileInput.files.length > 0 : false;

    if (topicChanged || nameChanged || descChanged || fileChanged) {
      submitButton.classList.remove("btn-disabled");
    } else {
      submitButton.classList.add("btn-disabled");
    }
  };

  const checkEmpty = () => {
    const topicEmpty = topicInput ? topicInput.value === "" : false;
    const nameEmpty = nameInput.value === "";
    const descEmpty = descTextarea.value === "";

    topicEmpty
      ? topicInput?.classList.add("input-error")
      : topicInput?.classList.remove("input-error");

    nameEmpty
      ? nameInput?.classList.add("input-error")
      : nameInput?.classList.remove("input-error");

    descEmpty
      ? descTextarea?.classList.add("textarea-error")
      : descTextarea?.classList.remove("textarea-error");

    if (descEmpty || nameEmpty || topicEmpty) {
      submitButton.classList.add("btn-disabled");
    }
  };

  // Проверяем при загрузке
  checkValidity();

  // Отслеживаем изменения
  [topicInput, nameInput, descTextarea].forEach((el) => {
    if (el) el.addEventListener("input", checkValidity);
    if (el) el.addEventListener("input", checkEmpty);
  });

  if (fileInput) fileInput.addEventListener("change", checkValidity);
}
