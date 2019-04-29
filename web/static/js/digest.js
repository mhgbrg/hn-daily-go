"use strict";

// Check if cached load. If it is, reload page to get new read status for clicked story.
const inputField = $("#refresh");
if (inputField.val() === "yes") {
  location.reload(true);
} else {
  inputField.val("yes");
}

// Set up clipboard and tooltips.
new ClipboardJS("#copy-device-id");
$("#device-id").tooltip();
$("#copy-device-id").tooltip({
  trigger: "click",
});

// Set up mark-as-read-buttons
$(".rank > form").submit(e => {
  e.preventDefault();
  const storyID = e.target.dataset.id;

  fetch(`/story/${storyID}/mark-as-read`, {
    method: "POST",
    body: JSON.stringify({ storyID }),
    headers: {
      "Content-Type": "application/json",
    },
  }).then(res => {
    if (!res.ok) {
      console.error("failed to mark story as read");
    }
  });

  $(`.title[data-id=${storyID}]`).addClass("is-read");
  $(e.target).parent().addClass("is-read");
});

// Mark story as read when clicked. This fixes a problem with stories not being marked as read
// when opened in a new tab.
$(".read-form").submit((e) => {
  const button = $(e.target).find(".title");
  const id = button.attr("data-id");
  const url = $(e.target).attr("action");
  const title = button.html();
  setTimeout(() => $(e.target).replaceWith(
    `<h3 class="title is-read slab mb-0" data-id="${id}"><a href="${url}">${title}</a></h3>`
  ));
});
