@extends('layouts.app')

@section('content')

    <div class="intro text-xl text-terminal font-pixel absolute left-28 top-64"></div>

    <div class="content hidden absolute left-28 top-64">
        <h1 class="text-2xl">Jurassic Park SNES Randomizer v0.3</h1>
        
        @if(session('error'))
            <div class="alert alert-error">
                {{ session('error') }}
            </div>
        @endif

        @if(session('success'))
            <div class="alert alert-success">
                {{ session('success') }}
            </div>
        @endif

        <form id="randomizerForm">
            @csrf
            
            <label class="block cursor-pointer w-fit">
            <!-- Das echte File-Input -->
                <input
                    type="file"
                    class="hidden"
                    id="rom"
                    name="rom"
                    accept=".sfc"
                    onchange="showFileName(event)"
                />

                <!-- Sichtbarer Terminal-Button -->
                <div
                    class="flex items-center px-4 py-2
                        border-4 border-terminal bg-black text-terminal font-pixel
                        hover:bg-green-900/20 transition-colors"
                >
                    <span id="fileLabel">[ Select ROM file ]</span>
                </div>
            </label>
            <p class="help-text mb-4">
                You need to provide a valid Jurassic Park v1.0 ROM file.<br />
                The ROM will be uploaded to the server for randomization and deleted afterwards.<br />
                The Jurassic Park Classic Game Collection on Steam includes a valid ROM file.
            </p>

            <div class="form-group mb-4">
                <label class="inline-block w-50" for="difficulty">Difficulty</label>
                <select class="border-4 border-terminal bg-black p-1" id="difficulty" name="difficulty" required>
                    <option value="0" selected>Easy - Only ID cards are randomized</option>
                    <option value="1" >Normal - ID cards and items per floor (softlocks possible)</option>
                    <option value="2">Hard - ID cards and items per building (softlocks possible)</option>
                </select>
            </div>

            <div class="form-group">
                <label for="seed" class="inline-block w-50">Seed (optional)</label>
                <input class="border-4 border-terminal bg-black p-1" type="number" min="0" step="1" id="seed" name="seed" placeholder="empty for random seed">
                <p class="help-text ml-50">Integer between 0 and 9223372036854775807</p>
            </div>

            <label class="ml-50 mt-4 flex items-center cursor-pointer select-none">
                <input
                    type="checkbox"
                    class="hidden peer"
                    id="startLocations" 
                    name="startLocations" 
                    value="1"
                />

                <span
                    class="w-5 h-5 flex items-center justify-center
                        border-4 border-terminal bg-black text-terminal
                        peer-checked:before:content-['X']
                        before:content-[''] before:text-xl"
                ></span>
                <span class="ml-2 text-terminal font-pixel">Random start location</span>
            </label>

            <label class="ml-50 mt-4 flex items-center cursor-pointer select-none">
                <input
                    type="checkbox"
                    class="hidden peer"
                    id="overworld" 
                    name="overworld" 
                    value="1"
                />

                <span
                    class="w-5 h-5 flex items-center justify-center
                        border-4 border-terminal bg-black text-terminal
                        peer-checked:before:content-['X']
                        before:content-[''] before:text-xl"
                ></span>
                <span class="ml-2 text-terminal font-pixel">Randomize overworld items</span>
            </label>

            <button type="submit" class="cursor-pointer btn border-4 p-4 my-4 border-terminal bg-black hover:bg-terminal hover:text-black" id="submitBtn">Start randomization</button>
        </form>

        <div id="status" class="border-4 border-terminal bg-black p-4 mt-4 hidden">
            <div class="" id="statusMessage"></div>
        </div>

        <div id="result" class="hidden">
            <div class="border-4 border-terminal bg-black p-4 mt-4" id="resultMessage">
                <strong>Randomization successful!</strong>
                <p>Seed: <span id="resultSeed"></span></p>                
            </div>

            <a class="cursor-pointer inline-block border-4 p-4 mt-4 border-terminal bg-black hover:bg-terminal hover:text-black" href="#" id="downloadLink" class="download-link">Download ROM</a>
        </div>
    </div>
@endsection

@push('scripts')
<script>
    const text = "PRESS ANY KEY TO SKIP...\n\nBIOS LOADED ...\n\nMESSAGE \nMAIN SYSTEM CHECK ... OK\n\nACTION\nINITIALIZING RANDOMIZER ...\n\nOK\n\nRANDOMIZER READY\n\nMESSAGE ENDS ...";
    const introElement = document.querySelector('.intro');

    document.onkeydown = function(e) {
        // Bei Tastendruck Intro überspringen
        showContent();
    };

    let index = 0;
    let speed = 35;
    function typeIntro() {
        if (index < text.length) {
            introElement.innerHTML += text.charAt(index) === '\n' ? '<br />' : text.charAt(index);
            index++;
            setTimeout(typeIntro, speed);
        } else {
            // Intro fertig, warten und text ausblenden
            setTimeout(showContent, 1000);
        }
    }

    function showContent() {
        document.querySelector('.intro').style.display = 'none';
        document.querySelector('.content').style.display = 'block';
    }

    typeIntro();
</script>
<script>
    document.getElementById('randomizerForm').addEventListener('submit', async function(e) {
        e.preventDefault();

        const formData = new FormData();
        const romFile = document.getElementById('rom').files[0];
        const submitBtn = document.getElementById('submitBtn');
        const statusDiv = document.getElementById('status');
        const statusMessage = document.getElementById('statusMessage');
        const resultDiv = document.getElementById('result');

        // Reset
        resultDiv.style.display = 'none';
        statusDiv.classList.remove('active');

        if (!romFile) {
            alert('Bitte wählen Sie eine ROM-Datei aus');
            return;
        }

        // Upload ROM
        submitBtn.disabled = true;
        submitBtn.textContent = 'Uploading...';
        statusDiv.classList.add('active');
        statusMessage.textContent = 'ROM upload and check...';

        try {
            const uploadFormData = new FormData();
            uploadFormData.append('rom', romFile);
            uploadFormData.append('_token', document.querySelector('meta[name="csrf-token"]').content);

            const uploadResponse = await fetch('{{ route("randomizer.upload") }}', {
                method: 'POST',
                body: uploadFormData,
            });

            const uploadResult = await uploadResponse.json();

            if (!uploadResult.success) {
                throw new Error(uploadResult.message || 'Upload failed');
            }

            // Start randomization
            submitBtn.textContent = 'Randomization in progress...';
            statusMessage.textContent = 'Randomization is being performed. This may take a few minutes...';

            const randomizeData = {
                romPath: uploadResult.path,
                difficulty: document.getElementById('difficulty').value,
                seed: document.getElementById('seed').value || '',
                startLocations: document.getElementById('startLocations').checked,
                overworld: document.getElementById('overworld').checked,
            };

            const randomizeResponse = await fetch('{{ route("randomizer.randomize") }}', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-CSRF-TOKEN': document.querySelector('meta[name="csrf-token"]').content,
                },
                body: JSON.stringify(randomizeData),
            });

            const randomizeResult = await randomizeResponse.json();

            if (!randomizeResult.success) {
                throw new Error(randomizeResult.message || 'Randomization failed');
            }

            // Erfolg
            statusDiv.classList.remove('active');
            resultDiv.style.display = 'block';
            document.getElementById('resultSeed').textContent = randomizeResult.seed;
            
            // Extrahiere den Pfad aus der Download-URL (Format: /download?path=...)
            let downloadPath = randomizeResult.download;
            if (downloadPath.includes('path=')) {
                downloadPath = downloadPath.split('path=')[1];
            }
            document.getElementById('downloadLink').href = '{{ route("randomizer.download") }}?path=' + encodeURIComponent(downloadPath);

            // Formular zurücksetzen
            document.getElementById('randomizerForm').reset();
            

        } catch (error) {
            statusDiv.classList.add('active');
            statusMessage.textContent = 'Fehler: ' + error.message;
            statusMessage.parentElement.className = 'alert alert-error';
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = 'Start randomization';
        }
    });
</script>
@endpush

