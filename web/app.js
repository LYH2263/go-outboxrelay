async function j(url, opts) {
  const r = await fetch(url, opts);
  return r.json();
}
async function refresh() {
  document.getElementById("list").textContent = JSON.stringify(await j("/api/events?limit=50"), null, 2);
  document.getElementById("stats").textContent = JSON.stringify(await j("/api/stats"), null, 2);
}
document.getElementById("btn-append").onclick = async () => {
  const topic = document.getElementById("topic").value;
  const url = document.getElementById("url").value;
  let payload;
  try { payload = JSON.parse(document.getElementById("payload").value); }
  catch { payload = document.getElementById("payload").value; }
  await j("/api/events/append", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ topic, payload, url }),
  });
  await refresh();
};
document.getElementById("btn-relay").onclick = async () => {
  const max = document.getElementById("max").value || 5;
  await j("/api/relay?max=" + max, { method: "POST" });
  await refresh();
};
refresh();
