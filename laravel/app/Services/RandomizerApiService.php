<?php

namespace App\Services;

use GuzzleHttp\Client;
use GuzzleHttp\Exception\GuzzleException;
use Illuminate\Http\UploadedFile;
use Illuminate\Support\Facades\Log;

class RandomizerApiService
{
    protected Client $client;
    protected string $baseUrl;

    public function __construct()
    {
        $this->baseUrl = config('services.randomizer.api_url');
        $this->client = new Client([
            'base_uri' => $this->baseUrl,
            'timeout' => 300, // 5 Minuten für Randomisierung
        ]);
    }

    /**
     * Lädt eine ROM-Datei zur Go-API hoch
     *
     * @param UploadedFile $file
     * @return array ['path' => string, 'name' => string]
     * @throws \Exception
     */
    public function uploadRom(UploadedFile $file): array
    {
        try {
            $response = $this->client->post('/upload', [
                'multipart' => [
                    [
                        'name' => 'rom',
                        'contents' => fopen($file->getRealPath(), 'r'),
                        'filename' => $file->getClientOriginalName(),
                    ],
                ],
            ]);

            $data = json_decode($response->getBody()->getContents(), true);

            if (!isset($data['path']) || !isset($data['name'])) {
                throw new \Exception('Ungültige Antwort von der API');
            }

            return $data;
        } catch (GuzzleException $e) {
            Log::error('ROM-Upload fehlgeschlagen', [
                'error' => $e->getMessage(),
                'response' => "",
            ]);

            $message = 'Fehler beim Upload der ROM-Datei';

            throw new \Exception($message, 0, $e);
        }
    }

    /**
     * Startet die Randomisierung mit den gegebenen Einstellungen
     *
     * @param array $settings
     * @return array ['download' => string, 'seed' => string]
     * @throws \Exception
     */
    public function randomize(array $settings): array
    {
        try {
            $response = $this->client->post('/randomize', [
                'json' => $settings,
                'headers' => [
                    'Content-Type' => 'application/json',
                ],
            ]);

            $data = json_decode($response->getBody()->getContents(), true);

            if (!isset($data['download']) || !isset($data['seed'])) {
                throw new \Exception('Ungültige Antwort von der API');
            }

            return $data;
        } catch (GuzzleException $e) {
            Log::error('Randomisierung fehlgeschlagen', [
                'error' => $e->getMessage(),
                'settings' => $settings,
                'response' => "",
            ]);

            $message = 'Fehler bei der Randomisierung';
            

            throw new \Exception($message, 0, $e);
        }
    }

    /**
     * Lädt eine randomisierte ROM-Datei herunter
     *
     * @param string $path
     * @return \Psr\Http\Message\StreamInterface
     * @throws \Exception
     */
    public function downloadRom(string $path)
    {
        try {
            $response = $this->client->get('/download', [
                'query' => ['path' => $path],
            ]);

            return $response->getBody();
        } catch (GuzzleException $e) {
            Log::error('ROM-Download fehlgeschlagen', [
                'error' => $e->getMessage(),
                'path' => $path,
            ]);

            throw new \Exception('Fehler beim Download der ROM-Datei', 0, $e);
        }
    }
}

