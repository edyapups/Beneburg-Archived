const schema: string = import.meta.env.VITE_SCHEMA ? import.meta.env.VITE_SCHEMA : 'http';
const host: string = import.meta.env.VITE_HOST ? import.meta.env.VITE_HOST : 'localhost:8080';
const url: string = `${schema}://${host}/api`;

export interface User {
    ID: string,
    CreatedAt: string,
    UpdatedAt: string,
    DeletedAt: string,
    telegram_id: string,
    username: string,
    name: string,
    age: string,
    sex: string,
    about: string,
    hobbies: string,
    work: string,
    education: string,
    cover_letter: string,
    contacts: string,
    is_bot: string,
    is_active: string,
}

export async function getUsers(): Promise<User[]> {
    const response = await fetch(url + '/users')
    return await response.json();
}

export async function getUser(id: string): Promise<User> {
    const response = await fetch(url + '/users/' + id)
    return await response.json();
}

export async function createUser(user: User): Promise<User> {
    const response = await fetch(url + '/users', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
    })
    return await response.json();
}

export async function updateUser(user: User): Promise<User> {
    const response = await fetch(url + '/users/' + user.ID, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
    })
    return await response.json();
}