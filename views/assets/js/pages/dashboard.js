(function () {
  const form = htmx.find("#file-upload-form");
  const input = htmx.find("#file-upload-input");
  const uploadBtn = htmx.find("#upload-button-topbar");
  const imageGallery = htmx.find("#dashboard-image-gallery");

  uploadBtn.addEventListener("click", (e) => {
    e.preventDefault();
    input.click();
  });

  // click anywhere in form
  form.addEventListener("click", (e) => {
    if (e.target.tagName !== "BUTTON") {
      input.click();
    }
  });

  // when file selected from "browse"
  input.addEventListener("change", () => {
    if (input.files.length > 0) {
      htmx.trigger(form, "submit");
    }
  });

  form.addEventListener("dragover", (e) => {
    e.preventDefault();
    form.classList.add("border-primary");
  });

  form.addEventListener("dragleave", () => {
    form.classList.remove("border-primary");
  });

  // Drop
  form.addEventListener("drop", (e) => {
    e.preventDefault();
    form.classList.remove("border-primary");

    if (e.dataTransfer.files.length > 0) {
      input.files = e.dataTransfer.files;
      htmx.trigger(form, "submit");
    }
  });

  // progress bar
  const progressContainer = document.getElementById(
    "upload-progress-container",
  );
  const progressBar = document.getElementById("upload-progress-bar");

  htmx.on(form, "htmx:xhr:loadstart", function () {
    progressContainer.classList.remove("hidden");
    progressBar.setAttribute("aria-valuenow", 0);
  });

  htmx.on(form, "htmx:xhr:progress", function (evt) {
    if (evt.detail.lengthComputable) {
      const percent = (evt.detail.loaded / evt.detail.total) * 100;
      progressBar.setAttribute("aria-valuenow", percent);
    }
  });

  htmx.on(form, "htmx:xhr:loadend", function () {
    setTimeout(() => {
      progressContainer.classList.add("hidden");
      progressBar.setAttribute("aria-valuenow", 0);
    }, 400);
  });

  // refresh gallery when image is uploaded
  htmx.on(form, "htmx:afterRequest", function (evt) {
    if (evt.target.id === "file-upload-form" && imageGallery) {
      imageGallery.dispatchEvent(new Event("refreshGallery"));
    }
  });
})();
