const previewImage = document.getElementById('preview-image');
const form = document.getElementById('generate-form');
const backgroundColorInput = document.getElementById('background-color');
const textColorInput = document.getElementById('text-color');
const randomizeBackgroundColorButton = document.getElementById('randomize-background-color');
const randomizeTextColorButton = document.getElementById('randomize-text-color');
const dowloadButton = document.getElementById('download-image');
const copyLink = document.getElementById('copy-url');
const copyImgTag = document.getElementById('copy-img-tag');
const copyLinkButton = document.getElementById('copy-url-btn');
const copyImgTagButton = document.getElementById('copy-img-btn');

const apiURL = 'https://localhost:8080';

let width = 1920;
let height = 1080;
let text = '';
let backgroundColor = 'dcdcdc';
let textColor = '000000';
let format = 'png';


// Functions
const generateImageURL = (options) => {
    options = {
        width: options.width || 1920,
        height: options.height || 1080,
        text: options.text || '',
        textColor: options.textColor || '000000',
        backgroundColor: options.backgroundColor || 'dcdcdc',
        format: options.format || 'jpeg',
    }

    let url = `${apiURL}/img/${options.width}x${options.height}?text-color=${options.textColor}&bg-color=${options.backgroundColor}&format=${options.format}`;

    if (options.text !== '') {
        url += `&text=${text}`;
    }

    return url;
}

const updatePreview = () => {
    previewImage.width = width;
    previewImage.height = height;

    const link = generateImageURL({
        width,
        height,
        text,
        textColor,
        backgroundColor,
        format,
    });

    previewImage.src = link;
    dowloadButton.href = link;
    copyLink.innerText = link;
    copyImgTag.innerText = `<img src="${link}" height="${height}" width="${width}" alt="placeholder-image" />`;

    if (text === '') {
        dowloadButton.download = `${width}x${height}.${format}`;
    } else {
        dowloadButton.download = `${text}.${format}`;
    }
}

const getRandomColor = () => {
    var letters = '0123456789ABCDEF';
    var color = '';
    for (var i = 0; i < 6; i++) {
      color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
}


// Event listeners
randomizeBackgroundColorButton.addEventListener('click', (e) => {
    e.preventDefault();

    backgroundColor = getRandomColor();
    backgroundColorInput.value = `#${backgroundColor}`;
});

randomizeTextColorButton.addEventListener('click', (e) => {
    e.preventDefault();

    textColor = getRandomColor();
    textColorInput.value = `#${textColor}`;
});

form.addEventListener('submit', (e) => {
    e.preventDefault();

    const formData = new FormData(form);
    width = formData.get('width');
    height = formData.get('height');
    text = formData.get('text');
    backgroundColor = formData.get('background-color').replace('#', '');
    textColor = formData.get('text-color').replace('#', '');
    format = formData.get('image-format');

    updatePreview();
});

copyLinkButton.addEventListener('click', (e) => {
    e.preventDefault();

    navigator.clipboard.writeText(copyLink.innerText);
    copyLinkButton.innerText = 'Copied!';
    setTimeout(() => {
        copyLinkButton.innerText = 'Copy';
    }, 2000);
});

copyImgTagButton.addEventListener('click', (e) => {
    e.preventDefault();

    navigator.clipboard.writeText(copyImgTag.innerText);
    copyImgTagButton.innerText = 'Copied!';
    setTimeout(() => {
        copyImgTagButton.innerText = 'Copy';
    }, 2000);
});

(() => {
    updatePreview();
})();
