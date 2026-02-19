import fs from "fs"
import path from "path"
import axios from "axios"
import sharp from "sharp"

// 1. configuration
const TARGET_SIZE = 128; // 128px is great for retina screens (64px logical)
const OUTPUT_DIR = './public/emojis'; // Where your Vue app looks for static files

// 2. The Master List (Paste the array I gave you earlier here)
const emojis = [
    { char: '😀', hex: '1f600', name: 'Grinning Face' },
    { char: '😃', hex: '1f603', name: 'Grinning Face with Big Eyes' },
    { char: '😄', hex: '1f604', name: 'Grinning Face with Smiling Eyes' },
    { char: '😁', hex: '1f601', name: 'Beaming Face' },
    { char: '😆', hex: '1f606', name: 'Grinning Squinting Face' },
    { char: '😅', hex: '1f605', name: 'Grinning Face with Sweat' },
    { char: '🤣', hex: '1f923', name: 'Rolling on the Floor Laughing' },
    { char: '😂', hex: '1f602', name: 'Face with Tears of Joy' },
    { char: '🙂', hex: '1f642', name: 'Slightly Smiling Face' },
    { char: '🙃', hex: '1f643', name: 'Upside-Down Face' },
    { char: '😉', hex: '1f609', name: 'Winking Face' },
    { char: '😊', hex: '1f60a', name: 'Smiling Face with Smiling Eyes' },
    { char: '😇', hex: '1f607', name: 'Smiling Face with Halo' },
    { char: '🥰', hex: '1f970', name: 'Smiling Face with Hearts' },
    { char: '😍', hex: '1f60d', name: 'Smiling Face with Heart-Eyes' },
    { char: '🤩', hex: '1f929', name: 'Star-Struck' },
    { char: '😘', hex: '1f618', name: 'Face Blowing a Kiss' },
    { char: '😗', hex: '1f617', name: 'Kissing Face' },
    { char: '😚', hex: '1f61a', name: 'Kissing Face with Closed Eyes' },
    { char: '😙', hex: '1f619', name: 'Kissing Face with Smiling Eyes' },
    { char: '😋', hex: '1f60b', name: 'Face Savoring Food' },
    { char: '😛', hex: '1f61b', name: 'Face with Tongue' },
    { char: '😜', hex: '1f61c', name: 'Winking Face with Tongue' },
    { char: '🤪', hex: '1f92a', name: 'Zany Face' },
    { char: '😝', hex: '1f61d', name: 'Squinting Face with Tongue' },
    { char: '🤑', hex: '1f911', name: 'Money-Mouth Face' },
    { char: '🤗', hex: '1f917', name: 'Hugging Face' },
    { char: '🤭', hex: '1f92d', name: 'Face with Hand Over Mouth' },
    { char: '🤫', hex: '1f92b', name: 'Shushing Face' },
    { char: '🤔', hex: '1f914', name: 'Thinking Face' },
    { char: '🤐', hex: '1f910', name: 'Zipper-Mouth Face' },
    { char: '🤨', hex: '1f928', name: 'Face with Raised Eyebrow' },
    { char: '😐', hex: '1f610', name: 'Neutral Face' },
    { char: '😑', hex: '1f611', name: 'Expressionless Face' },
    { char: '😶', hex: '1f636', name: 'Face Without Mouth' },
    { char: '😏', hex: '1f60f', name: 'Smirking Face' },
    { char: '😒', hex: '1f612', name: 'Unamused Face' },
    { char: '🙄', hex: '1f644', name: 'Face with Rolling Eyes' },
    { char: '😬', hex: '1f62c', name: 'Grimacing Face' },
    { char: '🤥', hex: '1f925', name: 'Lying Face' },
    { char: '😌', hex: '1f60c', name: 'Relieved Face' },
    { char: '😔', hex: '1f614', name: 'Pensive Face' },
    { char: '😪', hex: '1f62a', name: 'Sleepy Face' },
    { char: '🤤', hex: '1f924', name: 'Drooling Face' },
    { char: '😴', hex: '1f634', name: 'Sleeping Face' },
    { char: '😷', hex: '1f637', name: 'Face with Medical Mask' },
    { char: '🤒', hex: '1f912', name: 'Face with Thermometer' },
    { char: '🤕', hex: '1f915', name: 'Face with Head-Bandage' },
    { char: '🤢', hex: '1f922', name: 'Nauseated Face' },
    { char: '🤮', hex: '1f92e', name: 'Face Vomiting' },
    { char: '🤧', hex: '1f927', name: 'Sneezing Face' },
    { char: '🥵', hex: '1f975', name: 'Hot Face' },
    { char: '🥶', hex: '1f976', name: 'Cold Face' },
    { char: '🥴', hex: '1f974', name: 'Woozy Face' },
    { char: '😵', hex: '1f635', name: 'Dizzy Face' },
    { char: '🤯', hex: '1f92f', name: 'Exploding Head' },
    { char: '🤠', hex: '1f920', name: 'Cowboy Hat Face' },
    { char: '🥳', hex: '1f973', name: 'Partying Face' },
    { char: '😎', hex: '1f60e', name: 'Smiling Face with Sunglasses' },
    { char: '🤓', hex: '1f913', name: 'Nerd Face' },
    { char: '🧐', hex: '1f9d0', name: 'Face with Monocle' },
    { char: '😕', hex: '1f615', name: 'Confused Face' },
    { char: '😟', hex: '1f61f', name: 'Worried Face' },
    { char: '🙁', hex: '1f641', name: 'Slightly Frowning Face' },
    { char: '😮', hex: '1f62e', name: 'Face with Open Mouth' },
    { char: '😯', hex: '1f62f', name: 'Hushed Face' },
    { char: '😲', hex: '1f632', name: 'Astonished Face' },
    { char: '😳', hex: '1f633', name: 'Flushed Face' },
    { char: '🥺', hex: '1f97a', name: 'Pleading Face' },
    { char: '😦', hex: '1f626', name: 'Frowning Face with Open Mouth' },
    { char: '😧', hex: '1f627', name: 'Anguished Face' },
    { char: '😨', hex: '1f628', name: 'Fearful Face' },
    { char: '😰', hex: '1f630', name: 'Anxious Face with Sweat' },
    { char: '😥', hex: '1f625', name: 'Sad but Relieved Face' },
    { char: '😢', hex: '1f622', name: 'Crying Face' },
    { char: '😭', hex: '1f62d', name: 'Loudly Crying Face' },
    { char: '😱', hex: '1f631', name: 'Face Screaming in Fear' },
    { char: '😖', hex: '1f616', name: 'Confounded Face' },
    { char: '😣', hex: '1f623', name: 'Persevering Face' },
    { char: '😞', hex: '1f61e', name: 'Disappointed Face' },
    { char: '😓', hex: '1f613', name: 'Downcast Face with Sweat' },
    { char: '😩', hex: '1f629', name: 'Weary Face' },
    { char: '😫', hex: '1f62b', name: 'Tired Face' },
    { char: '😤', hex: '1f624', name: 'Face with Steam From Nose' },
    { char: '😡', hex: '1f621', name: 'Pouting Face' },
    { char: '😠', hex: '1f620', name: 'Angry Face' },
    { char: '😈', hex: '1f608', name: 'Smiling Face with Horns' },
    { char: '👿', hex: '1f47f', name: 'Angry Face with Horns' },
    { char: '💀', hex: '1f480', name: 'Skull' },
    { char: '💩', hex: '1f4a9', name: 'Pile of Poo' },
    { char: '🤡', hex: '1f921', name: 'Clown Face' },
    { char: '👻', hex: '1f47b', name: 'Ghost' },
    { char: '👽', hex: '1f47d', name: 'Alien' },
    { char: '🤖', hex: '1f916', name: 'Robot' },
    { char: '🙈', hex: '1f648', name: 'See-No-Evil Monkey' },
    { char: '🙉', hex: '1f649', name: 'Hear-No-Evil Monkey' },
    { char: '🙊', hex: '1f64a', name: 'Speak-No-Evil Monkey' },

    // --- ANIMALS & NATURE (Distinct Animations) ---
    { char: '🐱', hex: '1f431', name: 'Cat Face' },
    { char: '🐮', hex: '1f42e', name: 'Cow Face' },
    { char: '🐸', hex: '1f438', name: 'Frog' },
    { char: '🦄', hex: '1f984', name: 'Unicorn' },
    { char: '🐙', hex: '1f419', name: 'Octopus' },
    { char: '🐢', hex: '1f422', name: 'Turtle' },
    { char: '🦕', hex: '1f995', name: 'Sauropod' },
    { char: '🦖', hex: '1f996', name: 'T-Rex' },
    { char: '🐉', hex: '1f409', name: 'Dragon' },
    { char: '🦦', hex: '1f9a6', name: 'Otter' },
    { char: '🔥', hex: '1f525', name: 'Fire' },
    { char: '💧', hex: '1f4a7', name: 'Droplet' },
    { char: '🌈', hex: '1f308', name: 'Rainbow' },
    { char: '🌞', hex: '1f31e', name: 'Sun with Face' },

    // --- OBJECTS (Good for Game Items) ---
    { char: '💣', hex: '1f4a3', name: 'Bomb' },
    { char: '🎈', hex: '1f388', name: 'Balloon' },
    { char: '🎉', hex: '1f389', name: 'Party Popper' },
    { char: '🏆', hex: '1f3c6', name: 'Trophy' },
    { char: '⚽', hex: '26bd',  name: 'Soccer Ball' },
    { char: '💡', hex: '1f4a1', name: 'Light Bulb' },
    { char: '💎', hex: '1f48e', name: 'Gem Stone' },
    { char: '🎁', hex: '1f381', name: 'Wrapped Gift' },
    { char: '🚀', hex: '1f680', name: 'Rocket' },
    { char: '🛸', hex: '1f6f8', name: 'Flying Saucer' },
];

// Ensure directory exists
if (!fs.existsSync(OUTPUT_DIR)){
    fs.mkdirSync(OUTPUT_DIR, { recursive: true });
}

async function downloadAndResize(emoji) {
    const url = `https://fonts.gstatic.com/s/e/notoemoji/latest/${emoji.hex}/512.webp`;
    const outputPath = path.join(OUTPUT_DIR, `${emoji.hex}.webp`);

    try {
        // A. Download the file as a buffer
        const response = await axios({ url, responseType: 'arraybuffer' });
        
        // B. Resize using Sharp (maintaining animation)
        await sharp(response.data, { animated: true })
            .resize(TARGET_SIZE)
            .webp({ effort: 6 }) // Higher effort = smaller file size
            .toFile(outputPath);

        console.log(`✅ Saved: ${emoji.char} (${emoji.hex})`);
    } catch (error) {
        console.error(`❌ Failed: ${emoji.char} - ${error.message}`);
    }
}

async function run() {
    console.log(`🚀 Starting download of ${emojis.length} emojis...`);
    
    // Process in chunks to avoid banning (or simple serial loop)
    for (const emoji of emojis) {
        await downloadAndResize(emoji);
    }
    
    console.log("✨ All done! Your game assets are ready.");
}

run();