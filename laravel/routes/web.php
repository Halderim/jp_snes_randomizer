<?php

use App\Http\Controllers\RandomizerController;
use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| Web Routes
|--------------------------------------------------------------------------
|
| Here is where you can register web routes for your application. These
| routes are loaded by the RouteServiceProvider and all of them will
| be assigned to the "web" middleware group. Make something great!
|
*/

Route::get('/', [RandomizerController::class, 'index'])->name('randomizer.index');
Route::post('/upload', [RandomizerController::class, 'upload'])->name('randomizer.upload');
Route::post('/randomize', [RandomizerController::class, 'randomize'])->name('randomizer.randomize');
Route::get('/download', [RandomizerController::class, 'download'])->name('randomizer.download');
