"use strict";

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

// Mark story as read when clicked. This fixes two problems:
// 1. Stories not being marked as read when opened in a new tab
// 2. Cache not being cleared when opening story in same tab and then using back button to navigate
//    back to hn-daily
$(".read-form").submit((e) => {
  const button = $(e.target).find(".title");
  const id = button.attr("data-id");
  const url = $(e.target).attr("action");
  const title = button.html();
  setTimeout(() => $(e.target).replaceWith(
    `<h3 class="title is-read slab mb-0" data-id="${id}"><a href="${url}">${title}</a></h3>`
  ));
});
