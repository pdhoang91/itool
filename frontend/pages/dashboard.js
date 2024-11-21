// frontend/pages/dashboard.js

import { useEffect, useState } from 'react';
import { fetchTasks } from '../services/api';
import Header from '../components/Header';
import Footer from '../components/Footer';
import TaskTable from '../components/TaskTable/TaskTable';
import styles from '../styles/Dashboard.module.css';

export default function Dashboard() {
    const [tasks, setTasks] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    useEffect(() => {
        const getTasks = async () => {
            setLoading(true);
            setError(null);
            try {
                const response = await fetchTasks();
                setTasks(response.data);
            } catch (error) {
                console.error('Error fetching tasks:', error);
                setError('Có lỗi xảy ra khi tải danh sách tác vụ.');
            } finally {
                setLoading(false);
            }
        };

        getTasks();
    }, []);

    return (
        <div className={styles.container}>
            <Header />
            <main className={styles.main}>
                <h1>Dashboard</h1>
                {loading && <p>Đang tải...</p>}
                {error && <p className="error">{error}</p>}
                {!loading && !error && <TaskTable tasks={tasks} />}
            </main>
            <Footer />
        </div>
    );
}
