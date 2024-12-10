import React from 'react'
import {createRoot} from 'react-dom/client'
import AppRoute from './AppRoute'

// CSS App Styling Override and extra.
import './main.css';

const container = document.getElementById('root')

const root = createRoot(container)

root.render(
    <React.StrictMode>
        <AppRoute/>
    </React.StrictMode>
)
