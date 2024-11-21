// frontend/pages/index.js

import Link from 'next/link';
import styles from '../styles/Home.module.css';

export default function Home() {
    return (
        <div className={styles.container}>
            <h1>AI Tools Dashboard</h1>
            <ul>
                <li><Link href="/upload">Upload & Process</Link></li>
                <li><Link href="/dashboard">Dashboard</Link></li>
            </ul>
        </div>
    );
}
