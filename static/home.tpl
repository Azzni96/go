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

      .auth-section {
        background: #fff;
        border-radius: 10px;
        box-shadow: 0 0 10px rgba(0,0,0,0.1);
        margin-bottom: 20px;
        overflow: hidden;
      }

      .auth-header {
        background: #007bff;
        color: white;
        padding: 15px 20px;
        display: flex;
        justify-content: space-between;
        align-items: center;
      }

      .auth-form {
        padding: 20px;
        display: none;
      }

      .auth-form.active {
        display: block;
      }

      .form-group {
        margin-bottom: 15px;
      }

      .form-label {
        display: block;
        margin-bottom: 5px;
        font-weight: bold;
        color: #333;
      }

      .form-input {
        width: 100%;
        padding: 12px;
        border: 2px solid #ddd;
        border-radius: 5px;
        font-size: 16px;
      }

      .form-input:focus {
        outline: none;
        border-color: #007bff;
      }

      .btn-auth {
        width: 100%;
        padding: 12px;
        background: #007bff;
        color: white;
        border: none;
        border-radius: 5px;
        font-size: 16px;
        font-weight: bold;
        cursor: pointer;
        margin-bottom: 10px;
      }

      .btn-auth:hover {
        background: #0056b3;
      }

      .btn-secondary {
        background: #6c757d;
        width: 100%;
        padding: 12px;
        color: white;
        border: none;
        border-radius: 5px;
        font-size: 14px;
        cursor: pointer;
      }

      .btn-secondary:hover {
        background: #545b62;
      }

      .user-info {
        display: flex;
        align-items: center;
        gap: 15px;
      }

      .user-email {
        font-size: 14px;
        opacity: 0.9;
      }

      .btn-logout {
        background: #dc3545;
        color: white;
        border: none;
        padding: 8px 16px;
        border-radius: 5px;
        font-size: 12px;
        cursor: pointer;
      }

      .btn-logout:hover {
        background: #c82333;
      }

      .error-message {
        color: #dc3545;
        font-size: 14px;
        margin-top: 5px;
        display: none;
      }

      .success-message {
        color: #28a745;
        font-size: 14px;
        margin-top: 5px;
        display: none;
      }

      .auth-links {
        text-align: center;
        margin-top: 15px;
      }

      .auth-link {
        color: #007bff;
        cursor: pointer;
        text-decoration: underline;
        font-size: 14px;
      }

      .auth-link:hover {
        color: #0056b3;
      }
    </style>
  </head>
  <body>
    <!-- Authentication Section -->
    <div id="authSection" class="auth-section" style="max-width: 600px; margin: 0 auto;">
      <div class="auth-header">
        <div>
          <span id="authTitle">Welcome! Please login or register</span>
        </div>
        <div id="userInfo" class="user-info" style="display: none;">
          <span id="userEmail" class="user-email"></span>
          <button id="logoutBtn" class="btn-logout">Logout</button>
        </div>
      </div>

      <!-- Login Form -->
      <div id="loginForm" class="auth-form">
        <div class="form-group">
          <label class="form-label">Username or Email:</label>
          <input type="text" id="loginUsername" class="form-input" placeholder="Enter username or email">
        </div>
        <div class="form-group">
          <label class="form-label">Password:</label>
          <input type="password" id="loginPassword" class="form-input" placeholder="Enter password">
        </div>
        <button id="loginBtn" class="btn-auth">Login</button>
        <div id="loginError" class="error-message"></div>
        <div class="auth-links">
          <span class="auth-link" onclick="showRegisterForm()">Don't have an account? Register here</span>
        </div>
      </div>

      <!-- Register Form -->
      <div id="registerForm" class="auth-form">
        <div class="form-group">
          <label class="form-label">Username:</label>
          <input type="text" id="registerUsername" class="form-input" placeholder="Choose username">
        </div>
        <div class="form-group">
          <label class="form-label">Email:</label>
          <input type="email" id="registerEmail" class="form-input" placeholder="Enter email">
        </div>
        <div class="form-group">
          <label class="form-label">Password:</label>
          <input type="password" id="registerPassword" class="form-input" placeholder="Choose password (min 6 chars)">
        </div>
        <button id="registerBtn" class="btn-auth">Register</button>
        <div id="registerError" class="error-message"></div>
        <div id="registerSuccess" class="success-message"></div>
        <div class="auth-links">
          <span class="auth-link" onclick="showLoginForm()">Already have an account? Login here</span>
        </div>
      </div>
    </div>

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
      let currentUser = null;
      let authToken = localStorage.getItem('authToken');
      
      // DOM elements
      const todoForm = document.getElementById('todoForm');
      const todoInput = document.getElementById('todoInput');
      const addBtn = document.getElementById('addBtn');
      const todoList = document.getElementById('todoList');
      const emptyState = document.getElementById('emptyState');
      
      // Load todos when page loads
      document.addEventListener('DOMContentLoaded', function() {
        initializeAuth();
        loadTodos();
      });

      // Initialize authentication
      function initializeAuth() {
        if (authToken) {
          // Verify token and get user info
          fetch('/auth/profile', {
            headers: {
              'Authorization': 'Bearer ' + authToken
            }
          })
          .then(response => {
            if (response.ok) {
              return response.json();
            } else {
              throw new Error('Invalid token');
            }
          })
          .then(data => {
            if (data.data) {
              currentUser = data.data;
              showAuthenticatedState();
            } else {
              logout();
            }
          })
          .catch(() => {
            logout();
          });
        } else {
          showLoginForm();
        }
      }

      // Show login form
      function showLoginForm() {
        document.getElementById('loginForm').classList.add('active');
        document.getElementById('registerForm').classList.remove('active');
        document.getElementById('authTitle').textContent = 'Login to Your Account';
        document.getElementById('userInfo').style.display = 'none';
      }

      // Show register form  
      function showRegisterForm() {
        document.getElementById('registerForm').classList.add('active');
        document.getElementById('loginForm').classList.remove('active');
        document.getElementById('authTitle').textContent = 'Create New Account';
        document.getElementById('userInfo').style.display = 'none';
      }

      // Show authenticated state
      function showAuthenticatedState() {
        document.getElementById('loginForm').classList.remove('active');
        document.getElementById('registerForm').classList.remove('active');
        document.getElementById('authTitle').textContent = 'Welcome back!';
        document.getElementById('userInfo').style.display = 'flex';
        document.getElementById('userEmail').textContent = currentUser.email;
      }

      // Login function
      function login(username, password) {
        return fetch('/auth/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            username: username,
            password: password
          })
        })
        .then(response => response.json())
        .then(data => {
          if (data.error === false && data.data) {
            authToken = data.data.token;
            currentUser = data.data.user;
            localStorage.setItem('authToken', authToken);
            showAuthenticatedState();
            loadTodos(); // Reload todos for this user
            return true;
          } else {
            throw new Error(data.message || 'Login failed');
          }
        });
      }

      // Register function
      function register(username, email, password) {
        return fetch('/auth/register', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            username: username,
            email: email,
            password: password
          })
        })
        .then(response => response.json())
        .then(data => {
          if (data.error === false && data.data) {
            authToken = data.data.token;
            currentUser = data.data.user;
            localStorage.setItem('authToken', authToken);
            showAuthenticatedState();
            loadTodos(); // Load todos for new user
            return true;
          } else {
            throw new Error(data.message || 'Registration failed');
          }
        });
      }

      // Logout function
      function logout() {
        authToken = null;
        currentUser = null;
        localStorage.removeItem('authToken');
        todos = [];
        showLoginForm();
        renderTodos();
      }
      
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
        const headers = {};
        if (authToken) {
          headers['Authorization'] = 'Bearer ' + authToken;
        }

        fetch('/todo', { headers })
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
        const headers = {
          'Content-Type': 'application/json',
        };
        if (authToken) {
          headers['Authorization'] = 'Bearer ' + authToken;
        }

        fetch('/todo', {
          method: 'POST',
          headers: headers,
          body: JSON.stringify({ title: title })
        })
        .then(response => response.json())
        .then(data => {
          if (data.todo_id && data.todo) {
            // Add the new todo to local array
            todos.push(data.todo);
            renderTodos();
            document.querySelector('.todo-input').value = ''; // Clear input
            console.log('Todo added successfully!');
          } else {
            console.error('Invalid response format:', data);
          }
        })
        .catch(error => {
          console.error('Error adding todo:', error);
          alert('Virhe lisättäessä todoa. Yritä uudelleen.');
        });
      }
      
      // Update todo
      function updateTodo(id, title) {
        const todo = todos.find(t => t.id === id);
        if (!todo) return;
        
        const headers = {
          'Content-Type': 'application/json',
        };
        if (authToken) {
          headers['Authorization'] = 'Bearer ' + authToken;
        }
        
        fetch(`/todo/${id}`, {
          method: 'PUT',
          headers: headers,
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
        
        const headers = {
          'Content-Type': 'application/json',
        };
        if (authToken) {
          headers['Authorization'] = 'Bearer ' + authToken;
        }
        
        fetch(`/todo/${id}`, {
          method: 'PUT',
          headers: headers,
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
        
        const headers = {};
        if (authToken) {
          headers['Authorization'] = 'Bearer ' + authToken;
        }
        
        fetch(`/todo/${id}`, {
          method: 'DELETE',
          headers: headers
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

      // Auth form event listeners
      document.getElementById('loginBtn').addEventListener('click', function() {
        const username = document.getElementById('loginUsername').value.trim();
        const password = document.getElementById('loginPassword').value;
        const errorEl = document.getElementById('loginError');
        
        if (!username || !password) {
          errorEl.textContent = 'Please fill in all fields';
          errorEl.style.display = 'block';
          return;
        }

        this.textContent = 'Logging in...';
        this.disabled = true;
        errorEl.style.display = 'none';

        login(username, password)
          .then(() => {
            // Success - form will be hidden by showAuthenticatedState()
          })
          .catch(error => {
            errorEl.textContent = error.message;
            errorEl.style.display = 'block';
          })
          .finally(() => {
            this.textContent = 'Login';
            this.disabled = false;
          });
      });

      document.getElementById('registerBtn').addEventListener('click', function() {
        const username = document.getElementById('registerUsername').value.trim();
        const email = document.getElementById('registerEmail').value.trim();
        const password = document.getElementById('registerPassword').value;
        const errorEl = document.getElementById('registerError');
        const successEl = document.getElementById('registerSuccess');
        
        if (!username || !email || !password) {
          errorEl.textContent = 'Please fill in all fields';
          errorEl.style.display = 'block';
          successEl.style.display = 'none';
          return;
        }

        if (password.length < 6) {
          errorEl.textContent = 'Password must be at least 6 characters';
          errorEl.style.display = 'block';
          successEl.style.display = 'none';
          return;
        }

        this.textContent = 'Registering...';
        this.disabled = true;
        errorEl.style.display = 'none';
        successEl.style.display = 'none';

        register(username, email, password)
          .then(() => {
            successEl.textContent = 'Registration successful! Welcome!';
            successEl.style.display = 'block';
          })
          .catch(error => {
            errorEl.textContent = error.message;
            errorEl.style.display = 'block';
          })
          .finally(() => {
            this.textContent = 'Register';
            this.disabled = false;
          });
      });

      document.getElementById('logoutBtn').addEventListener('click', function() {
        if (confirm('Are you sure you want to logout?')) {
          logout();
        }
      });

      // Allow Enter key to submit forms
      document.getElementById('loginPassword').addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
          document.getElementById('loginBtn').click();
        }
      });

      document.getElementById('registerPassword').addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
          document.getElementById('registerBtn').click();
        }
      });
    </script>
  </body>
</html>