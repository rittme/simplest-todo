const populateLists = async (entries) => {
  const todo = document.getElementById('to-do-list');
  const done = document.getElementById('done-list');

  const todoFragment = document.createDocumentFragment();
  entries.todo.forEach(({ label, id }) => {
    const item = document.createElement('li');
    const text = document.createTextNode(label);

    const del = document.createElement('a');
    del.classList.add('action');
    del.href = `/entry/${id}`;
    del.dataset.method = 'DELETE';
    del.appendChild(document.createTextNode('🗑'));

    const done = document.createElement('a');
    done.classList.add('action');
    done.href = `/entry/${id}`;
    done.dataset.method = 'PUT';
    done.appendChild(document.createTextNode('✅'));

    item.appendChild(text);
    item.appendChild(del);
    item.appendChild(done);
    todoFragment.appendChild(item);
  });

  const doneFragment = document.createDocumentFragment();
  entries.done.forEach(({ label, id }) => {
    const item = document.createElement('li');
    const text = document.createTextNode(label);

    const del = document.createElement('a');
    del.classList.add('action');
    del.href = `/entry/${id}`;
    del.dataset.method = 'DELETE';
    del.appendChild(document.createTextNode('🗑'));

    const done = document.createElement('a');
    done.classList.add('action');
    done.href = `/entry/${id}`;
    done.dataset.method = 'PUT';
    done.appendChild(document.createTextNode('↩️'));

    item.appendChild(text);
    item.appendChild(del);
    item.appendChild(done);
    doneFragment.appendChild(item);
  });

  todo.innerHTML = '';
  todo.appendChild(todoFragment);
  done.innerHTML = '';
  done.appendChild(doneFragment);
};

//////////////////////////////////////////////////////////////////////////

const myUrl = '';

const loadContent = async () => {
  const response = await fetch(`${myUrl}/entry`);
  return response.json();
};

const getEntries = async (content) => {
  const { entries } = content;
  return entries.reduce(
    (acc, entry, i) => {
      const list = entry.done ? acc.done : acc.todo;
      list.push({ label: entry.label, id: i });
      return acc;
    },
    { done: [], todo: [] }
  );
};

const handleAction = async (event) => {
  if (!event.target.classList.contains('action')) {
    return true;
  }
  event.preventDefault();

  const path = event.target.href;
  const method = event.target.dataset.method;
  
  const response = await fetch(path, {
    method,
  });
  const content = await response.json();

  getEntries(content).then(populateLists);
};

loadContent().then((content) => getEntries(content).then(populateLists));
document.getElementById('lists').addEventListener('click', handleAction);

//////////////////////////////////////////////////////////////////////////

let currentRandomId;

const getRandom = async () => {
  const response = await fetch(`${myUrl}/random`);
  return response.json();
};

const populateRandom = async () => {
  try {
    const { Entry: entry, ID: id } = await getRandom();

    document.getElementById('hero').innerText = entry.label;
    currentRandomId = id;
  } catch (_) {
    document.getElementById('hero').innerText = '';
    currentRandomId = undefined;
  }
};

document.getElementById('hero-new').addEventListener('click', () => populateRandom());
document.getElementById('hero-done').addEventListener('click', async function (){
  const id = currentRandomId;
  
  const response = await fetch(`/entry/${id}`, {
    method: 'PUT',
  });
  const content = await response.json();

  getEntries(content).then(populateLists);
  populateRandom();
});

populateRandom();

//////////////////////////////////////////////////////////////////////////

document.getElementById('add').addEventListener('keyup', async function (event) {
  event.preventDefault();
  if (event.code == 'Enter' && this.value.length > 0) {
    const response = await fetch('/entry', {
      method: 'POST',
      body: this.value,
    });

    const content = await response.json();
    getEntries(content).then(populateLists);
    this.value = '';
  }
})