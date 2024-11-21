// frontend/components/TaskTable/TaskRow.js

import styles from '../../styles/TaskRow.module.css';

export default function TaskRow({ task }) {
    return (
        <tr className={styles.row}>
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
    );
}
