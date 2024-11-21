// frontend/components/Upload/UploadSection.js

import TextToVoice from './TextToVoice';
import VoiceToText from './VoiceToText';
import RemoveBackground from './RemoveBackground';
import SpeechRecognition from './SpeechRecognition';
import FaceRecognition from './FaceRecognition';
import OCR from './OCR';
import Translation from './Translation';
import styles from '../../styles/UploadSection.module.css';

export default function UploadSection({ result, setResult }) {
    return (
        <div className={styles.uploadSection}>
            <TextToVoice setResult={setResult} />
            <VoiceToText setResult={setResult} />
            <RemoveBackground setResult={setResult} />
            <SpeechRecognition setResult={setResult} />
            <FaceRecognition setResult={setResult} />
            <OCR setResult={setResult} />
            <Translation setResult={setResult} />
            {result && (
                <div className={styles.result}>
                    <h3>Result:</h3>
                    <pre>{JSON.stringify(result, null, 2)}</pre>
                </div>
            )}
        </div>
    );
}
