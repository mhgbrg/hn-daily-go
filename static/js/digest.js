'use strict';

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
