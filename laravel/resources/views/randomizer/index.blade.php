@extends('layouts.app')

@section('content')
    <h1>🦖 Jurassic Park SNES Randomizer</h1>
    <p class="subtitle">Randomisiere dein Jurassic Park SNES Spiel</p>

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

        <div class="form-group">
            <label for="rom">ROM-Datei hochladen (.smc oder .sfc)</label>
            <input type="file" id="rom" name="rom" accept=".smc,.sfc" required>
            <p class="help-text">Die ROM muss eine erweiterte 4MB Version von Jurassic Park USA v1.0 sein</p>
        </div>

        <div class="form-group">
            <label for="difficulty">Schwierigkeitsgrad</label>
            <select id="difficulty" name="difficulty" required>
                <option value="0">Easy - Nur ID-Karten werden randomisiert</option>
                <option value="1" selected>Normal - ID-Karten und Items pro Etage</option>
                <option value="2">Hard - ID-Karten und Items pro Gebäude</option>
            </select>
        </div>

        <div class="form-group">
            <label for="seed">Seed (optional)</label>
            <input type="text" id="seed" name="seed" placeholder="Leer lassen für zufälligen Seed">
            <p class="help-text">Geben Sie einen Seed ein, um die gleiche Randomisierung zu reproduzieren</p>
        </div>

        <div class="form-group">
            <div class="checkbox-group">
                <input type="checkbox" id="startLocations" name="startLocations" value="1">
                <label for="startLocations">Zufällige Startposition</label>
            </div>
        </div>

        <div class="form-group">
            <div class="checkbox-group">
                <input type="checkbox" id="overworld" name="overworld" value="1" checked>
                <label for="overworld">Overworld-Items randomisieren</label>
            </div>
        </div>

        <button type="submit" class="btn" id="submitBtn">Randomisierung starten</button>
    </form>

    <div id="status" class="status">
        <div class="alert alert-info" id="statusMessage"></div>
    </div>

    <div id="result" style="display: none;">
        <div class="alert alert-success">
            <strong>Randomisierung erfolgreich!</strong>
            <p>Seed: <span id="resultSeed"></span></p>
            <a href="#" id="downloadLink" class="download-link">ROM herunterladen</a>
        </div>
    </div>
@endsection

@push('scripts')
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
        submitBtn.textContent = 'ROM wird hochgeladen...';
        statusDiv.classList.add('active');
        statusMessage.textContent = 'ROM wird hochgeladen und geprüft...';

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
                throw new Error(uploadResult.message || 'Upload fehlgeschlagen');
            }

            // Randomisierung starten
            submitBtn.textContent = 'Randomisierung läuft...';
            statusMessage.textContent = 'Randomisierung wird durchgeführt. Dies kann einige Minuten dauern...';

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
                throw new Error(randomizeResult.message || 'Randomisierung fehlgeschlagen');
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
            document.getElementById('difficulty').value = '1';
            document.getElementById('overworld').checked = true;

        } catch (error) {
            statusDiv.classList.add('active');
            statusMessage.textContent = 'Fehler: ' + error.message;
            statusMessage.parentElement.className = 'alert alert-error';
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = 'Randomisierung starten';
        }
    });
</script>
@endpush

