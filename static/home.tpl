<!DOCTYPE html>
<html lang="en">
  <head>
    <title>Simple Todo App</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
      * {
        box-sizing: border-box;
        margin: 0;
        padding: 0;
      }
      
      body {
        font-family: Arial, sans-serif;
        background-color: #f4f4f4;
        padding: 20px;
      }
      
      .container {
        max-width: 600px;
        margin: 0 auto;
        background: white;
        border-radius: 10px;
        box-shadow: 0 0 10px rgba(0,0,0,0.1);
        overflow: hidden;
      }
      
      .header {
        background: #333;
        color: white;
        padding: 20px;
        text-align: center;
      }
      
      .header h1 {
        font-size: 24px;
      }
      
      .todo-form {
        padding: 20px;
        border-bottom: 1px solid #eee;
      }
      
      .input-group {
        display: flex;
        gap: 10px;
      }
      
      .todo-input {
        flex: 1;
        padding: 12px;
        border: 2px solid #ddd;
        border-radius: 5px;
        font-size: 16px;
      }
      
      .todo-input:focus {
        outline: none;
        border-color: #007bff;
      }
      
      .todo-input.error {
        border-color: #dc3545;
      }
      
      .btn {
        padding: 12px 20px;
        border: none;
        border-radius: 5px;
        cursor: pointer;
        font-size: 16px;
        font-weight: bold;
      }
      
      .btn-add {
        background: #28a745;
        color: white;
      }
      
      .btn-edit {
        background: #ffc107;
        color: black;
      }
      
      .btn-small {
        padding: 5px 10px;
        font-size: 12px;
        margin: 0 2px;
      }
      
      .btn-success {
        background: #28a745;
        color: white;
      }
      
      .btn-danger {
        background: #dc3545;
        color: white;
      }
      
      .todo-list {
        list-style: none;
      }
      
      .todo-item {
        padding: 15px 20px;
        border-bottom: 1px solid #eee;
        display: flex;
        align-items: center;
        justify-content: space-between;
        cursor: pointer;
        transition: background-color 0.3s;
      }
      
      .todo-item:hover {
        background-color: #f8f9fa;
      }
      
      .todo-item.completed {
        background-color: #e9ecef;
        color: #6c757d;
      }
      
      .todo-text {
        flex: 1;
        margin-left: 10px;
      }
      
      .todo-text.completed {
        text-decoration: line-through;
      }
      
      .todo-checkbox {
        width: 18px;
        height: 18px;
        cursor: pointer;
      }
      
      .todo-actions {
        display: flex;
        gap: 5px;
      }
      
      .empty-state {
        padding: 40px 20px;
        text-align: center;
        color: #6c757d;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="header">
        <h1>Simple Todo App</h1>
      </div>
      
      <div class="todo-form">
        <form id="todoForm">
          <div class="input-group">
            <input 
              type="text" 
              id="todoInput" 
              class="todo-input" 
              placeholder="Add your todo..."
              autocomplete="off"
            >
            <button type="submit" id="addBtn" class="btn btn-add">Add</button>
          </div>
        </form>
      </div>
      
      <ul id="todoList" class="todo-list">
        <!-- Todos will be dynamically added here -->
      </ul>
      
      <div id="emptyState" class="empty-state" style="display: none;">
        <p>No todos yet. Add one above!</p>
      </div>
    </div>
    <script>
      let todos = [];
      let editingTodo = null;
      
      // DOM elements
      const todoForm = document.getElementById('todoForm');
      const todoInput = document.getElementById('todoInput');
      const addBtn = document.getElementById('addBtn');
      const todoList = document.getElementById('todoList');
      const emptyState = document.getElementById('emptyState');
      
      // Load todos when page loads
      document.addEventListener('DOMContentLoaded', function() {
        loadTodos();
      });
      
      // Form submission
      todoForm.addEventListener('submit', function(e) {
        e.preventDefault();
        const title = todoInput.value.trim();
        
        if (title === '') {
          todoInput.classList.add('error');
          setTimeout(() => {
            todoInput.classList.remove('error');
          }, 2000);
          return;
        }
        
        if (editingTodo) {
          updateTodo(editingTodo.id, title);
        } else {
          addTodo(title);
        }
        
        todoInput.value = '';
        resetEditMode();
      });
      
      // Load todos from server
      function loadTodos() {
        fetch('/todo')
          .then(response => response.json())
          .then(data => {
            todos = data.data || [];
            renderTodos();
          })
          .catch(error => {
            console.error('Error loading todos:', error);
          });
      }
      
      // Add new todo
      function addTodo(title) {
        fetch('/todo', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ title: title })
        })
        .then(response => response.json())
        .then(data => {
          if (data.todo_id) {
            todos.push({
              id: data.todo_id,
              title: title,
              completed: false
            });
            renderTodos();
          }
        })
        .catch(error => {
          console.error('Error adding todo:', error);
        });
      }
      
      // Update todo
      function updateTodo(id, title) {
        const todo = todos.find(t => t.id === id);
        if (!todo) return;
        
        fetch(`/todo/${id}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            id: id,
            title: title,
            completed: todo.completed
          })
        })
        .then(response => {
          if (response.status === 200) {
            todo.title = title;
            renderTodos();
          }
        })
        .catch(error => {
          console.error('Error updating todo:', error);
        });
      }
      
      // Toggle todo completion
      function toggleTodo(id) {
        const todoIndex = todos.findIndex(t => t.id === id);
        if (todoIndex === -1) return;
        
        const todo = todos[todoIndex];
        const newCompleted = !todo.completed;
        
        fetch(`/todo/${id}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            id: id,
            title: todo.title,
            completed: newCompleted
          })
        })
        .then(response => {
          if (response.status === 200) {
            todos[todoIndex].completed = newCompleted;
            renderTodos();
          }
        })
        .catch(error => {
          console.error('Error toggling todo:', error);
        });
      }
      
      // Edit todo
      function editTodo(id) {
        const todo = todos.find(t => t.id === id);
        if (!todo) return;
        
        editingTodo = todo;
        todoInput.value = todo.title;
        addBtn.textContent = 'Update';
        addBtn.className = 'btn btn-edit';
        todoInput.focus();
      }
      
      // Delete todo
      function deleteTodo(id) {
        if (!confirm('Are you sure you want to delete this todo?')) {
          return;
        }
        
        fetch(`/todo/${id}`, {
          method: 'DELETE'
        })
        .then(response => {
          if (response.status === 200) {
            todos = todos.filter(t => t.id !== id);
            renderTodos();
            resetEditMode();
          }
        })
        .catch(error => {
          console.error('Error deleting todo:', error);
        });
      }
      
      // Reset edit mode
      function resetEditMode() {
        editingTodo = null;
        addBtn.textContent = 'Add';
        addBtn.className = 'btn btn-add';
      }
      
      // Render todos
      function renderTodos() {
        todoList.innerHTML = '';
        
        if (todos.length === 0) {
          emptyState.style.display = 'block';
          return;
        }
        
        emptyState.style.display = 'none';
        
        todos.forEach(todo => {
          const li = document.createElement('li');
          li.className = `todo-item ${todo.completed ? 'completed' : ''}`;
          
          li.innerHTML = `
            <input 
              type="checkbox" 
              class="todo-checkbox" 
              ${todo.completed ? 'checked' : ''}
              onclick="toggleTodo('${todo.id}')"
            >
            <span class="todo-text ${todo.completed ? 'completed' : ''}">${todo.title}</span>
            <div class="todo-actions">
              <button class="btn btn-small btn-success" onclick="editTodo('${todo.id}')">Edit</button>
              <button class="btn btn-small btn-danger" onclick="deleteTodo('${todo.id}')">Delete</button>
            </div>
          `;
          
          todoList.appendChild(li);
        });
      }
    </script>
  </body>
</html>