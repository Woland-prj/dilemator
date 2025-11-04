const form = document.getElementById("editor-form");
if (form) {
  const topicInput = form.querySelector('input[name="Topic"]');
  const nameInput = form.querySelector('input[name="Name"]');
  const descTextarea = form.querySelector('textarea[name="Value"]');
  const submitButton = form.querySelector('button[type="submit"]');

  if (!submitButton) return;

  const checkValidity = () => {
    const topicValid = topicInput ? topicInput.value.trim() !== "" : true;
    const nameValid = nameInput && nameInput.value.trim() !== "";
    const descValid = descTextarea && descTextarea.value.trim() !== "";

    if (topicValid && nameValid && descValid) {
      submitButton.classList.remove("btn-disabled");
    } else {
      submitButton.classList.add("btn-disabled");
    }
  };

  // Проверяем при загрузке
  checkValidity();

  // Отслеживаем изменения
  [topicInput, nameInput, descTextarea].forEach((el) => {
    if (el) el.addEventListener("input", checkValidity);
  });
}
