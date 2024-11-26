// frontend/pages/upload.js

import { useState } from 'react';
import Header from '../components/Header';
import Footer from '../components/Footer';
import UploadSection from '../components/Upload/UploadSection';
import styles from '../styles/Upload.module.css';
import TextToSpeech from '../components/Upload/TextToVoice';

export default function Upload() {
    const [result, setResult] = useState(null);

    return (
        <div className={styles.container}>
            <Header />
            <main className={styles.main}>
                <h1>Upload & Process</h1>
                <UploadSection result={result} setResult={setResult} />
                <TextToSpeech />
            </main>
            <Footer />
        </div>
    );
}
