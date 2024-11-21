// frontend/components/TaskTable/TaskTable.js

import TaskRow from './TaskRow';
import styles from '../../styles/TaskTable.module.css';

export default function TaskTable({ tasks }) {
    return (
        <table className={styles.table}>
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
                    <TaskRow key={task.id} task={task} />
                ))}
            </tbody>
        </table>
    );
}
