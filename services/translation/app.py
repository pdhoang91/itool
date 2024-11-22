from flask import Flask, request, jsonify
from googletrans import Translator

app = Flask(__name__)
translator = Translator()

@app.route('/translate', methods=['POST'])
def translate():
    data = request.get_json()
    if not data or 'text' not in data or 'dest_lang' not in data:
        return jsonify({"error": "Invalid input"}), 400

    text = data['text']
    dest_lang = data['dest_lang']

    try:
        # Thực hiện dịch văn bản
        translated = translator.translate(text, dest=dest_lang).text
        response = {"translated_text": translated}
    except Exception as e:
        response = {"error": str(e)}

    return jsonify(response), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5007)
