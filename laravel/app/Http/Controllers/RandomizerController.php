<?php

namespace App\Http\Controllers;

use App\Http\Requests\RandomizeRequest;
use App\Services\RandomizerApiService;
use Illuminate\Http\Request;
use Illuminate\Http\Response;
use Illuminate\Support\Facades\Log;

class RandomizerController extends Controller
{
    protected RandomizerApiService $apiService;

    public function __construct(RandomizerApiService $apiService)
    {
        $this->apiService = $apiService;
    }

    /**
     * Zeigt das Hauptformular
     */
    public function index()
    {
        return view('randomizer.index');
    }

    /**
     * Verarbeitet den ROM-Upload
     */
    public function upload(Request $request)
    {
        //dd($request);
        $request->validate([
            'rom' => ['required', 'file', 'max:4096'], // Max 4MB
        ], [
            'rom.required' => 'Bitte wählen Sie eine ROM-Datei aus',
            'rom.file' => 'Die hochgeladene Datei ist ungültig',
            'rom.max' => 'Die Datei darf maximal 4MB groß sein',
        ]);

        try {
            $result = $this->apiService->uploadRom($request->file('rom'));

            return response()->json([
                'success' => true,
                'path' => $result['path'],
                'name' => $result['name'],
            ]);
        } catch (\Exception $e) {
            Log::error('Upload-Fehler im Controller', [
                'error' => $e->getMessage(),
            ]);

            return response()->json([
                'success' => false,
                'message' => $e->getMessage(),
            ], 400);
        }
    }

    /**
     * Startet die Randomisierung
     */
    public function randomize(RandomizeRequest $request)
    {
        try {
            $settings = [
                'romPath' => $request->input('romPath'),
                'difficulty' => (int) $request->input('difficulty'),
                'seed' => $request->input('seed', ''),
                'startLocations' => $request->boolean('startLocations', false),
                'overworld' => $request->boolean('overworld', false),
            ];

            $result = $this->apiService->randomize($settings);

            return response()->json([
                'success' => true,
                'download' => $result['download'],
                'seed' => $result['seed'],
            ]);
        } catch (\Exception $e) {
            Log::error('Randomisierungs-Fehler im Controller', [
                'error' => $e->getMessage(),
            ]);

            return response()->json([
                'success' => false,
                'message' => $e->getMessage(),
            ], 500);
        }
    }

    /**
     * Stellt den ROM-Download bereit
     */
    public function download(Request $request)
    {
        $request->validate([
            'path' => ['required', 'string'],
        ]);

        try {
            $path = $request->input('path');
            
            // Extrahiere den Pfad aus der URL, falls eine vollständige URL übergeben wurde
            if (strpos($path, '/download?path=') !== false) {
                $path = urldecode(substr($path, strpos($path, 'path=') + 5));
            }
            
            $stream = $this->apiService->downloadRom($path);

            $filename = basename($path);

            return response()->streamDownload(function () use ($stream) {
                echo $stream->getContents();
            }, $filename, [
                'Content-Type' => 'application/octet-stream',
            ]);
        } catch (\Exception $e) {
            Log::error('Download-Fehler im Controller', [
                'error' => $e->getMessage(),
                'path' => $request->input('path'),
            ]);

            return redirect()->route('randomizer.index')
                ->with('error', 'Fehler beim Download: ' . $e->getMessage());
        }
    }
}

