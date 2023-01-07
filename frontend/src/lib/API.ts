import {getToken} from "./Token";

export declare interface User {
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
    const response = await fetch('/users', {
        headers: {
            'Authorization': `Bearer ${getToken()}`,
        }
    })
    return await response.json();
}

export async function getUser(id: string): Promise<User> {
    const response = await fetch('/users/' + id, {
        headers: {
            'Authorization': `Bearer ${getToken()}`,
        }
    })
    return await response.json();
}

export async function getMe(): Promise<User> {
    const response = await fetch('/getMe', {
        headers: {
            'Authorization': `Bearer ${getToken()}`,
        }
    })
    return await response.json();
}

export async function createUser(user: User): Promise<User> {
    const response = await fetch('/users', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${getToken()}`,
        },
        body: JSON.stringify(user)
    })
    return await response.json();
}

export async function updateUser(user: User): Promise<User> {
    const response = await fetch('/users/' + user.ID, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${getToken()}`,
        },
        body: JSON.stringify(user)
    })
    return await response.json();
}