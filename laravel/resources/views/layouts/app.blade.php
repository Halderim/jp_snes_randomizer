<!DOCTYPE html>
<html lang="de">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="csrf-token" content="{{ csrf_token() }}">
    @vite(['resources/css/app.css', 'resources/js/app.js'])
    <title>Jurassic Park SNES Randomizer</title>

    <link rel="preconnect" href="https://fonts.bunny.net">
    <link href="https://fonts.bunny.net/css?family=figtree:400,600&display=swap" rel="stylesheet" />
    <link href="https://fonts.bunny.net/css?family=pixelify-sans:400,500,600,700" rel="stylesheet" />
    
</head>
<body class="bg-black text-terminal font-pixel">
    <div class="container h-[1126px] w-[1281px] bg-no-repeat" style="background-image: url('/img/terminal.webp');">
        @yield('content')
    </div>

    @stack('scripts')
</body>
</html>

