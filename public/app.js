/* global fetch */
const packListDiv = document.getElementById('packList');
const packMsgEl   = document.getElementById('packMsg');
const calcMsgEl   = document.getElementById('calcMsg');
const resultTbl   = document.getElementById('resultTbl');
const resultBody  = resultTbl.querySelector('tbody');

document.getElementById('addPack').addEventListener('click', () => addPackRow());
document.getElementById('savePacks').addEventListener('click', savePacks);
document.getElementById('calcBtn').addEventListener('click', calc);

// ---------- INITIAL LOAD ----------
(async () => {
  try {
    const res  = await fetch('/packs');
    const data = res.ok ? await res.json() : null;
    showPackSizes(data?.length ? data : undefined);
  } catch {
    // If the backend is down use defaults
    showPackSizes();
  }
})();

function showPackSizes(arr = [250, 500, 1000, 2000, 5000]) {
  packListDiv.innerHTML = '';
  arr.forEach(addPackRow);
}

function addPackRow(value = '') {
  packListDiv.insertAdjacentHTML(
    'beforeend',
    `<div class="flex pack-row">
       <input type="number" min="1" value="${value}" />
       <button class="btn remove-row" title="Remove pack size">✖</button>
     </div>`
  );
  packListDiv.lastElementChild
    .querySelector('.remove-row')
    .addEventListener('click', e => e.currentTarget.parentElement.remove());
}

// ---------- SAVE PACK SIZES ----------
async function savePacks() {
  const sizes = [...packListDiv.querySelectorAll('input')]
    .map(i => Number(i.value))
    .filter(v => v > 0);

  if (!sizes.length) {
    return setMsg(packMsgEl, 'Enter at least one size');
  }

  try {
    const res = await fetch('/packs', {
      method : 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body   : JSON.stringify({ packSizes: sizes }),
    });

    if (!res.ok) throw new Error(await res.text());
    setMsg(packMsgEl, 'Saved!');
  } catch (err) {
    setMsg(packMsgEl, err.message || 'Failed to save');
  }
}

// ---------- CALCULATE ORDER ----------
async function calc() {
  const itemsInput = document.getElementById('items');
  const items      = Number(itemsInput.value);

  if (!items) {
    return setMsg(calcMsgEl, 'Enter a positive number');
  }

  try {
    const res = await fetch('/api/calculate', {
      method : 'POST',
      headers: { 'Content-Type': 'application/json' },
      body   : JSON.stringify({ items }),
    });
    if (!res.ok) throw new Error(await res.text());

    const { result } = await res.json();
    renderResult(result);
    setMsg(calcMsgEl, '');
  } catch (err) {
    setMsg(calcMsgEl, err.message || 'Calculation failed');
  }
}

function renderResult(result) {
  resultBody.innerHTML = '';
  Object.entries(result)
    .sort(([a], [b]) => b - a) // numeric desc on pack size
    .forEach(([pack, qty]) =>
      resultBody.insertAdjacentHTML(
        'beforeend',
        `<tr><td>${pack}</td><td>${qty}</td></tr>`
      )
    );
  resultTbl.classList.remove('hidden');
}

function setMsg(el, text) {
  el.textContent = text;
}
