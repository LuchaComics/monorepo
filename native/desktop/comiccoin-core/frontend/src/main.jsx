import React from 'react'
import {createRoot} from 'react-dom/client'
import './style.css'
import AppRoute from './AppRoute'

const container = document.getElementById('root')

const root = createRoot(container)

root.render(
    <React.StrictMode>
        <AppRoute/>
    </React.StrictMode>
)
