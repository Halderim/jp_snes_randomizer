<?php

namespace App\Http\Requests;

use Illuminate\Foundation\Http\FormRequest;

class RandomizeRequest extends FormRequest
{
    /**
     * Determine if the user is authorized to make this request.
     */
    public function authorize(): bool
    {
        return true;
    }

    /**
     * Get the validation rules that apply to the request.
     *
     * @return array<string, \Illuminate\Contracts\Validation\ValidationRule|array<mixed>|string>
     */
    public function rules(): array
    {
        return [
            'romPath' => ['required', 'string'],
            'difficulty' => ['required', 'integer', 'in:0,1,2'],
            'seed' => ['nullable', 'string'],
            'startLocations' => ['boolean'],
            'overworld' => ['boolean'],
        ];
    }

    /**
     * Get custom messages for validator errors.
     *
     * @return array<string, string>
     */
    public function messages(): array
    {
        return [
            'romPath.required' => 'ROM-Pfad ist erforderlich',
            'difficulty.required' => 'Schwierigkeitsgrad ist erforderlich',
            'difficulty.in' => 'Ungültiger Schwierigkeitsgrad. Muss 0, 1 oder 2 sein.',
        ];
    }

    /**
     * Prepare the data for validation.
     */
    protected function prepareForValidation(): void
    {
        $this->merge([
            'startLocations' => $this->has('startLocations'),
            'overworld' => $this->has('overworld'),
        ]);
    }
}

