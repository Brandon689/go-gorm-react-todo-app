import React, { useState, useEffect } from 'react';

function TodoApp() {
  const [todos, setTodos] = useState([]);
  const [todoText, setTodoText] = useState('');
  const [todoStatus, setTodoStatus] = useState('');

  useEffect(() => {
    fetchTodos();
  }, []);

  const fetchTodos = () => {
    fetch('http://localhost:8000/todos')
      .then(response => {
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.json();
      })
      .then(todosData => setTodos(todosData))
      .catch(error => console.error('Error fetching Todos:', error));
  };

  const addTodo = () => {
    fetch('http://localhost:8000/todos', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        text: todoText,
        status: todoStatus,
      }),
    })
      .then(response => response.json())
      .then(() => {
        // Clear input fields
        setTodoText('');
        setTodoStatus('');

        // Fetch and display updated Todos
        fetchTodos();
      })
      .catch(error => console.error('Error adding Todo:', error));
  };

  const removeTodo = async (todoID) => {
    //const encodedID = encodeURIComponent(todoID);

    try {
      const response = await fetch(`http://localhost:8000/todos/${todoID}`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        console.error('Error removing Todo:', response.statusText);
        return;
      }

      // Fetch and display updated Todos
      fetchTodos();
    } catch (error) {
      console.error('Error removing Todo:', error);
    }
  };

  return (
    <div>
      <h1>Todo App</h1>

      {/* Form to add a new Todo */}
      <form>
        <label htmlFor="todoText">Todo Text:</label>
        <input type="text" id="todoText" value={todoText} onChange={(e) => setTodoText(e.target.value)} required />
        <label htmlFor="todoStatus">Todo Status:</label>
        <input type="text" id="todoStatus" value={todoStatus} onChange={(e) => setTodoStatus(e.target.value)} required />
        <button type="button" onClick={addTodo}>Add Todo</button>
      </form>

      {/* Display Todos */}
      <ul>
        {todos.map(todo => (
          <li key={todo.id}>
            {`${todo.text} - ${todo.status}`}
            <button type="button" onClick={() => removeTodo(todo.id)}>Remove</button>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default TodoApp;
