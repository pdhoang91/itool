// frontend/pages/index.js

import Header from '../components/Header';
import Footer from '../components/Footer';
import styles from '../styles/Home.module.css';
import Link from 'next/link';

export default function Home() {
    return (
        <div className={styles.container}>
            <Header />
            <main className={styles.main}>
                <h1>AI Tools Dashboard</h1>
                <div className={styles.links}>
                    <Link href="/upload" className={styles.link}>Upload & Process</Link>
                    <Link href="/dashboard" className={styles.link}>Dashboard</Link>
                </div>
            </main>
            <Footer />
        </div>
    );
}
