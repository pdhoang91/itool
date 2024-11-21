// // frontend/pages/dashboard.js
import { useEffect, useState } from 'react';
import axios from 'axios';

export default function Dashboard() {
    const [tasks, setTasks] = useState([]);

    useEffect(() => {
        fetchTasks();
    }, []);

    const fetchTasks = async () => {
        try {
            // Gọi endpoint /tasks để lấy danh sách tất cả các task
            const response = await axios.get('http://localhost:81/tasks');
            setTasks(response.data);
        } catch (error) {
            console.error('Error fetching tasks:', error);
        }
    };

    return (
        <div>
            <h1>Dashboard</h1>
            <table>
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Service</th>
                        <th>Status</th>
                        <th>Created At</th>
                        <th>Updated At</th>
                        <th>Output</th>
                    </tr>
                </thead>
                <tbody>
                    {tasks.map(task => (
                        <tr key={task.id}>
                            <td>{task.id}</td>
                            <td>{task.service_name}</td>
                            <td>{task.status}</td>
                            <td>{new Date(task.created_at).toLocaleString()}</td>
                            <td>{new Date(task.updated_at).toLocaleString()}</td>
                            <td>
                                {task.output_data ? (
                                    <pre>{JSON.stringify(task.output_data, null, 2)}</pre>
                                ) : (
                                    'N/A'
                                )}
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}
