window.addEventListener("load", (_) => {
  try {
    const evtSource = new EventSource("/server/sent/event/browser/reload");
    evtSource.onmessage = function (_) {
      window.location.reload();
    };
  } catch (_) {
    window.location.reload();
  }
});
